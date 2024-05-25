package acceptance

import (
	"fmt"
	"os"
	"slices"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

type UserResource struct {
	Name           string
	Password       string
	Tags           []string
	MaxConnections string
	MaxChannels    string
}

func (u *UserResource) RequiredCreate(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
		password = "%s"
	}`, data.ResourceType, data.ResourceLabel, u.Name, u.Password)
}

func (u *UserResource) RequiredUpdate(data TestData) string {
	u.Password = data.RandomString()
	return u.RequiredCreate(data)
}

func (u *UserResource) OptionalCreate(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
		password = "%s"
		tags = %s
		max_connections = %s
		max_channels = %s
	}`, data.ResourceType, data.ResourceLabel, u.Name, u.Password, data.BuildArrayString(u.Tags), u.MaxConnections, u.MaxChannels)
}

func (u *UserResource) OptionalUpdateTags(data TestData) string {
	u.Tags = []string{"console"}
	return u.OptionalCreate(data)
}

func (u *UserResource) OptionalUpdateLimits(data TestData) string {
	u.MaxConnections = data.RandomIntegerString()
	u.MaxChannels = data.RandomIntegerString()
	return u.OptionalCreate(data)
}

func (u *UserResource) OptionalRemove(data TestData) string {
	u.Tags = []string{}
	u.MaxConnections = ""
	u.MaxChannels = ""
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
		password = "%s"
	}`, data.ResourceType, data.ResourceLabel, u.Name, u.Password)
}

func (u *UserResource) LoginCreate(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
		password = "%s"
		tags = %s
	}`, data.ResourceType, data.ResourceLabel, u.Name, u.Password, data.BuildArrayString(u.Tags))
}

func (u *UserResource) LoginUpdate(data TestData) string {
	u.Password = data.RandomString()
	return u.LoginCreate(data)
}

func (u *UserResource) ErrorConvertingCreate(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
		password = "%s"
		max_connections = "%s"
		max_channels = "%s"
	}`, data.ResourceType, data.ResourceLabel, u.Name, u.Password, u.MaxConnections, u.MaxChannels)
}

func (u *UserResource) ErrorConvertingUpdate(data TestData, connections string, channels string) string {
	u.MaxConnections = connections
	u.MaxChannels = channels
	return u.ErrorConvertingCreate(data)
}

func (u UserResource) ExistsInRabbitMQ() error {

	rmqc := TestAcc.Provider.Meta().(*rabbithole.Client)
	myUser, err := rmqc.GetUser(u.Name)
	if err != nil {
		return fmt.Errorf("error retrieving user '%s': %#v", u.Name, err)
	}
	if myUser.Name != u.Name {
		return fmt.Errorf("user name is not equal. Actual: '%s' Expected: %s", myUser.Name, u.Name)
	}
	if len(myUser.PasswordHash) <= 0 {
		return fmt.Errorf("user password is empty")
	}
	if len(myUser.Tags) != len(u.Tags) {
		return fmt.Errorf("user tags number is not equal. Actual: '%d' Expected: %d", len(myUser.Tags), len(u.Tags))
	} else {
		for i := 0; i < len(u.Tags); i++ {
			if !slices.Contains(myUser.Tags, myUser.Tags[i]) {
				return fmt.Errorf("user tags '%s' is not contained. Actual: '%#v' Expected: %#v", myUser.Tags, myUser.Tags, u.Tags[i])
			}
		}
	}

	myUserLimits, err := rmqc.GetUserLimits(u.Name)
	if err != nil {
		return fmt.Errorf("error retrieving user limit '%s': %#v", u.Name, err)
	}
	if len(myUserLimits) == 0 {
		if u.MaxConnections != "" && u.MaxChannels != "" {
			return fmt.Errorf("user limit is not empty: %#v", myUserLimits)
		}
	} else {
		if myUserLimits[0].Value["max-connections"] == 0 && u.MaxConnections == "" {
			//It's OK (specific case).
		} else if strconv.Itoa(myUserLimits[0].Value["max-connections"]) != u.MaxConnections {
			return fmt.Errorf("user limit 'max-connections' is not equal. Actual: '%d' Expected: %s", myUserLimits[0].Value["max-connections"], u.MaxConnections)
		}
		if myUserLimits[0].Value["max-channels"] == 0 && u.MaxChannels == "" {
			//It's OK (specific case).
		} else if strconv.Itoa(myUserLimits[0].Value["max-channels"]) != u.MaxChannels {
			return fmt.Errorf("user limit 'max-channels' is not equal. Actual: '%d' Expected: %s", myUserLimits[0].Value["max-channels"], u.MaxChannels)
		}
	}

	return nil
}

func (u UserResource) CheckLoginInRabbitMQ() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := rabbithole.NewClient(os.Getenv("RABBITMQ_ENDPOINT"), u.Name, u.Password)
		if err != nil {
			return fmt.Errorf("could not create RabbitMQ client: %#v", err)
		}

		_, err = client.Whoami()
		if err != nil {
			return fmt.Errorf("could not call whoami with username '%s': %#v", u.Name, err)
		}
		return nil
	}
}

func (u UserResource) CheckDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rmqc := TestAcc.Provider.Meta().(*rabbithole.Client)
		user, err := rmqc.GetUser(u.Name)
		if err != nil && err.(rabbithole.ErrorResponse).StatusCode != 404 {
			return fmt.Errorf("error retrieving user '%s': %#v", u.Name, err)
		}

		if user != nil {
			return fmt.Errorf("user still exists: %s", u.Name)
		}

		return nil
	}
}
