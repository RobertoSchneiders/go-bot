package bot

import (
	"errors"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMessageReceived(t *testing.T) {
	Convey("Given a new message in the channel", t, func() {
		commands = make(map[string]*customCommand)
		conn := &ircConnectionMock{}

		Convey("When the command is not registered", func() {
			conn = &ircConnectionMock{}

			Convey("It should not post to the channel", func() {

				messageReceived("#go-bot", "!not_a_cmd", "user", conn)

				So(conn.Messages, ShouldBeEmpty)
			})

		})

		Convey("The command can return an error", func() {
			conn = &ircConnectionMock{}

			Convey("it sould send the message with the error to the channel", func() {
				cmdError := errors.New("error")
				RegisterCommand("cmd", "", "",
					func(c *Cmd) (string, error) {
						return "", cmdError
					})

				messageReceived("#go-bot", "!cmd", "user", conn)

				So(conn.Channel, ShouldEqual, "#go-bot")
				So(conn.Messages, ShouldResemble,
					[]string{fmt.Sprintf(errorExecutingCommand, "cmd", cmdError.Error())})
			})
		})

		Convey("When the command is valid and registered", func() {
			conn = &ircConnectionMock{}

			commands = make(map[string]*customCommand)
			expectedMsg := "msg"
			cmd := "cmd"
			cmdDescription := "Command description"
			cmdExampleArgs := "arg1 arg2"

			RegisterCommand(cmd, cmdDescription, cmdExampleArgs,
				func(c *Cmd) (string, error) {
					return expectedMsg, nil
				})

			Convey("If it is called in the channel, reply on the channel", func() {
				messageReceived("#go-bot", "!cmd", "user", conn)

				So(conn.Channel, ShouldEqual, "#go-bot")
				So(conn.Messages, ShouldResemble, []string{expectedMsg})
			})

			Convey("If it is a private message, reply to the user", func() {
				conn.Nick = "go-bot"
				messageReceived("go-bot", "!cmd", "sender-nick", conn)

				So(conn.Channel, ShouldEqual, "sender-nick")
			})

			Convey("When the command is help", func() {

				Convey("Display the available commands in the channel", func() {
					messageReceived("#go-bot", "!help", "user", conn)

					So(conn.Channel, ShouldEqual, "#go-bot")
					So(conn.Messages, ShouldResemble, []string{
						fmt.Sprintf(helpAboutCommand, CmdPrefix),
						fmt.Sprintf(availableCommands, "cmd"),
					})
				})

				Convey("If the command exists send a message to the channel", func() {

					messageReceived("#go-bot", "!help cmd", "user", conn)

					So(conn.Channel, ShouldEqual, "#go-bot")
					So(conn.Messages, ShouldResemble, []string{
						fmt.Sprintf(helpDescripton, cmdDescription),
						fmt.Sprintf(helpUsage, CmdPrefix, cmd, cmdExampleArgs),
					})

				})

				Convey("If the command does not exists, display the generic help", func() {
					messageReceived("#go-bot", "!help not_a_command", "user", conn)

					So(conn.Channel, ShouldEqual, "#go-bot")
					So(conn.Messages, ShouldResemble, []string{
						fmt.Sprintf(helpAboutCommand, CmdPrefix),
						fmt.Sprintf(availableCommands, "cmd"),
					})
				})
			})

		})
	})

}
