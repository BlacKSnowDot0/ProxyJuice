package Core

import (
	"ProxyJuice/Utility"
	"fmt"
	"github.com/logrusorgru/aurora/v4"
	"os"
	"strings"
)

func SuccessChannelHandler() {
	for data := range SuccessChannel {
		var str strings.Builder
		str.WriteString(strings.ToLower(string(data.Type)))
		str.WriteString("://")

		if data.Username != "" && data.Password != "" {
			str.WriteString(data.Username)
			str.WriteString(":")
			str.WriteString(data.Password)
			str.WriteString("@")
		}

		str.WriteString(data.Address)
		str.WriteString(":")
		str.WriteString(Utility.FastIntToStr(data.Port))

		Utility.AppendToFile(OutputFile, str.String()+"\n")

		LogChannel <- &JuicyLog{Message: fmt.Sprintf("Proxy Detected: %s", str.String()), Level: SUCCESS, Sender: "Core.SuccessChannelHandler"}
	}
}

func LogChannelHandler() {
	exit := false
	for data := range LogChannel {
		var msg strings.Builder
		msg.WriteString("[")
		switch data.Level {
		case SUCCESS:
			msg.WriteString(aurora.Green("SUCCESS").String())

		case OK:
			msg.WriteString(aurora.Green("OK").String())

		case WARN:
			msg.WriteString(aurora.Yellow("WARN").String())

		case BAD:
			msg.WriteString(aurora.Red("BAD").String())

		case DIE:
			msg.WriteString(aurora.Red("DIE").String())
			exit = true
		}

		msg.WriteString("] [")
		msg.WriteString(aurora.White(data.Sender).String())
		msg.WriteString("] ")

		msg.WriteString(data.Message)

		fmt.Println(msg.String())

		if exit {
			os.Exit(1)
		}
	}
}
