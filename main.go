package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	Token   string `split_words:"true" required:"true"`
	GuildID string `split_words:"true" required:"true"`
}

type discordHandler struct {
	config Config
}

type discordServerIDs struct {
	ParentIDDiscordServer string `json:"id_discord_server"`
	ParentIDChannelText   string `json:"id_text_channel"`
	ParentIDVoiceText     string `json:"id_voice_channel"`
}

func main() {
	var config Config
	config.Token = ""
	config.GuildID = ""
	/*err := envconfig.Process("bot", &config)
	  if Token != nil {
	  	fmt.Printf("Error with configuration: %s\n", err.Error())
	  	os.Exit(1)
	  }*/

	discord, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Printf("Error creating Discord session: %s\n", err.Error())
		os.Exit(1)
	}

	// Register callbacks
	dh := &discordHandler{}
	discord.AddHandler(dh.ready)
	discord.AddHandler(dh.command)

	err = discord.Open()
	if err != nil {
		fmt.Printf("Error opening Discord connection: %s\n", err.Error())
		os.Exit(1)
	}

	commandDeleteChannel := &discordgo.ApplicationCommand{
		Name:        "delete-channel",
		Description: "Delete text and voice channels",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "first-channel",
				Description: "Pick the channel that you want to delete",
				Type:        discordgo.ApplicationCommandOptionChannel,
				Required:    false,
			},
			{
				Name:        "second-channel",
				Description: "Pick the channel that you want to delete",
				Type:        discordgo.ApplicationCommandOptionChannel,
				Required:    false,
			},
			{
				Name:        "third-channel",
				Description: "Pick the channel that you want to delete",
				Type:        discordgo.ApplicationCommandOptionChannel,
				Required:    false,
			},
			{
				Name:        "fourth-channel",
				Description: "Pick the channel that you want to delete",
				Type:        discordgo.ApplicationCommandOptionChannel,
				Required:    false,
			},
			{
				Name:        "fifth-channel",
				Description: "Pick the channel that you want to delete",
				Type:        discordgo.ApplicationCommandOptionChannel,
				Required:    false,
			},
			{
				Name:        "sixth-channel",
				Description: "Pick the channel that you want to delete",
				Type:        discordgo.ApplicationCommandOptionChannel,
				Required:    false,
			},
			{
				Name:        "seventh-channel",
				Description: "Pick the channel that you want to delete",
				Type:        discordgo.ApplicationCommandOptionChannel,
				Required:    false,
			},
			{
				Name:        "eight-channel",
				Description: "Pick the channel that you want to delete",
				Type:        discordgo.ApplicationCommandOptionChannel,
				Required:    false,
			},
			{
				Name:        "ninth-channel",
				Description: "Pick the channel that you want to delete",
				Type:        discordgo.ApplicationCommandOptionChannel,
				Required:    false,
			},
			{
				Name:        "tenth-channel",
				Description: "Pick the channel that you want to delete",
				Type:        discordgo.ApplicationCommandOptionChannel,
				Required:    false,
			},
		},
	}

	commandArchive := &discordgo.ApplicationCommand{
		Name:        "archive",
		Description: "Archive the chat for its participants",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "role",
				Description: "The role to archive",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
		},
	}

	commandCreateTeam := &discordgo.ApplicationCommand{
		Name:        "create-team",
		Description: "Create a new team",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "team",
				Description: "The name of the team",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "first-member",
				Description: "Team member",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "second-member",
				Description: "Team member",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "third-member",
				Description: "Team member",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "fourth-member",
				Description: "Team member",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "fifth-member",
				Description: "Team member",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "sixth-member",
				Description: "Team member",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "seventh-ember",
				Description: "Team member",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "eighth-ember",
				Description: "Team member",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "custom-color",
				Description: "pick your custom color (HEX type of color!!!!)",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "text-category",
				Description: "set a category where your text chat is going to be located",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "voice-category",
				Description: "set a category where your voice chat is going to be located",
				Required:    false,
			},
		},
	}

	commandRating := &discordgo.ApplicationCommand{
		Name:        "rating",
		Description: "Get rating of the player",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "player",
				Description: "The osu id or username of the player",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
		},
	}

	cmdArchive, err := discord.ApplicationCommandCreate(discord.State.User.ID, config.GuildID, commandArchive)
	if err != nil {
		fmt.Printf("Error adding command: %s\n", err.Error())
	}

	cmdCreateTeam, err := discord.ApplicationCommandCreate(discord.State.User.ID, config.GuildID, commandCreateTeam)
	if err != nil {
		fmt.Printf("Error adding command: %s\n", err.Error())
	}
	cmdCommandRating, err := discord.ApplicationCommandCreate(discord.State.User.ID, config.GuildID, commandRating)
	if err != nil {
		fmt.Printf("Error adding command: %s\n", err.Error())
	}
	cmdCommandDeleteChannel, err := discord.ApplicationCommandCreate(discord.State.User.ID, config.GuildID, commandDeleteChannel)
	if err != nil {
		fmt.Printf("Error adding command: %s\n", err.Error())
	}

	// Block until we get ctrl-c
	fmt.Println("Bot running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-sc

	// Clean up
	fmt.Println("Bot exiting")
	err = discord.ApplicationCommandDelete(discord.State.User.ID, config.GuildID, cmdArchive.ID)
	if err != nil {
		fmt.Printf("Error removing command: %s\n", err.Error())
	}
	err = discord.ApplicationCommandDelete(discord.State.User.ID, config.GuildID, cmdCreateTeam.ID)
	if err != nil {
		fmt.Printf("Error removing command: %s\n", err.Error())
	}
	err = discord.ApplicationCommandDelete(discord.State.User.ID, config.GuildID, cmdCommandRating.ID)
	if err != nil {
		fmt.Printf("Error removing command: %s\n", err.Error())
	}
	err = discord.ApplicationCommandDelete(discord.State.User.ID, config.GuildID, cmdCommandDeleteChannel.ID)
	if err != nil {
		fmt.Printf("Error removing command: %s\n", err.Error())
	}
	discord.Close()
}

func (dh *discordHandler) ready(s *discordgo.Session, m *discordgo.Ready) {
	s.UpdateListeningStatus("Listening")
}

func (dh *discordHandler) command(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		switch i.ApplicationCommandData().Name {
		case "archive":
			handleArchive(s, i)
		case "create-team":
			handleCreateTeam(s, i)
		case "rating":
			handleRating(s, i)
		case "delete-channel":
			handleDeletingChannels(s, i)
		default:
			Respond(s, i, "WHO ARE YOU? DIDN'T READ, LOL")
		}
	} else if i.Type == discordgo.InteractionMessageComponent {
		RespondForThinking(s, i)
		member, err := s.GuildMember(i.GuildID, i.Member.User.ID)
		if err != nil {
			fmt.Println("error retrieving member,", err)
			return
		}
		if !hasRole(member.Roles, "Секретарь ЦК импрува", s, i.GuildID) {
			Respond(s, i, "DIDN'T ASK")
			return
		}
		if strings.HasPrefix(i.MessageComponentData().CustomID, "delete_channel") {
			deleteSelectedChannel(s, i)
		} else if i.MessageComponentData().CustomID == "cancel_delete_channel" {
			Respond(s, i, "Channel deletion canceled.")
			err = s.ChannelMessageDelete(i.ChannelID, i.Message.ID)
			if err != nil {
				fmt.Println("Error deleting message:", err)
			}
		} else if i.MessageComponentData().CustomID == "select_channel" {
			channelName := i.MessageComponentData().Values[0]
			Respond(s, i, fmt.Sprintf("You selected channel: %s", channelName))
		} else {
			Respond(s, i, "Waiting for approval")
		}
		Respond(s, i, "  ") // TODO: тут что-то придумать, поскольку после применения вертится по кд
	} else {
		Respond(s, i, "This interaction type is not supported.")
	}
}

func deleteSelectedChannel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID
	channelNames := strings.Split(strings.TrimPrefix(customID, "delete_channel:"), ",")
	for index := range channelNames {
		channelNames[index] = strings.TrimSpace(channelNames[index])
	}
	channels, err := s.GuildChannels(i.GuildID)
	if err != nil {
		fmt.Println("Error retrieving channels:", err)
		Respond(s, i, "Error retrieving channels.")
		return
	}
	for _, channel := range channels {
		// Проверяем, есть ли имя канала в списке channelNames
		for _, channelName := range channelNames {
			if channel.Name == channelName {
				if channel.Type == discordgo.ChannelTypeGuildVoice {
					_, err := s.ChannelDelete(channel.ID)
					if err != nil {
						fmt.Println("Error deleting channel:", err)
						Respond(s, i, "Error deleting channel: "+err.Error())
						return
					}
					err = s.ChannelMessageDelete(i.ChannelID, i.Message.ID)
					if err != nil {
						fmt.Println("Error deleting message:", err)
					}
				}
			}
		}
	}
	username := i.Member.User.Username
	currentDate := time.Now().Format("January 2, 2006")
	archiveMessage := fmt.Sprintf("Chat archived since %s by `%s`. Voice chats were removed.", currentDate, username)
	_, err = s.ChannelMessageSend(i.ChannelID, archiveMessage)
	if err != nil {
		fmt.Println("Error sending archival message:", err)
	}
	Respond(s, i, "Selected channels have been processed.")
}

func RespondForThinking(s *discordgo.Session, i *discordgo.InteractionCreate) {
	typeRespond := discordgo.InteractionResponseDeferredChannelMessageWithSource
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: typeRespond,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Error sending message to discord: %v", err)
		return
	}
}

func RespondForThinkingVisible(s *discordgo.Session, i *discordgo.InteractionCreate) {
	typeRespond := discordgo.InteractionResponseDeferredChannelMessageWithSource
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: typeRespond,
	})
	if err != nil {
		log.Printf("Error sending message to discord: %v", err)
		return
	}
}

func Respond(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	var err error
	_, err = s.FollowupMessageCreate(i.Interaction,
		true,
		&discordgo.WebhookParams{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral})
	if err != nil {
		log.Printf("Error sending message to discord: %v", err)
		return
	}
}

func CheckHexColor(s string) bool {
	// Regular expression to match HEX color codes
	re := regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`) // ^ - start of string, # - start of hex, 0-9 - min and max num, A-F - min max sym,
	return re.MatchString(s)                      // a-f min max sym but smoll, 6 - six symbols (# doesn't count), $ - end of string
}

func handleCreateTeam(s *discordgo.Session, i *discordgo.InteractionCreate) {
	RespondForThinking(s, i)
	member, err := s.GuildMember(i.GuildID, i.Member.User.ID)
	if err != nil {
		fmt.Println("error retrieving member,", err)
		return
	}

	if hasRole(member.Roles, "Секретарь ЦК импрува", s, i.GuildID) {
		//var userIDs []string
		//user := options[1].UserValue(s)

		//locals
		options := i.ApplicationCommandData().Options
		roleName := options[0].StringValue()
		color := getRandomColor()
		ParentIDDiscordServer := "0"
		ParentIDChannelText := "0"
		ParentIDVoiceText := "0"

		// json
		fileBytes, err := ioutil.ReadFile("./discord_settings.json")
		if err != nil {
			panic(err)
		}

		var ids discordServerIDs

		err = json.Unmarshal(fileBytes, &ids)
		if err != nil {
			panic(err)
		}

		// fmt.Println("ID of Discord Server:", ids.ParentIDDiscordServer)
		// fmt.Println("ID of Text Channel:", ids.ParentIDChannelText)
		// fmt.Println("ID of Voice Channel:", ids.ParentIDVoiceText)

		// fmt.Println("[orig] ParentIDChannelText - ", ParentIDChannelText)
		// fmt.Println("[orig] ParentIDVoiceText - ", ParentIDVoiceText)

		if ParentIDDiscordServer == ids.ParentIDDiscordServer {
			fmt.Println("yup, all good")
			ParentIDChannelText = ids.ParentIDChannelText
			ParentIDVoiceText = ids.ParentIDVoiceText
			fmt.Println("[yup] ParentIDChannelText - ", ParentIDChannelText)
			fmt.Println("[yup] ParentIDVoiceText - ", ParentIDVoiceText)
		} else {
			fmt.Println("[bad] ID of Discord Server is different from ID from discord_settings.json file. Text category and voice category are going to be at the firsts categories")

			channels, err := s.GuildChannels(i.GuildID)
			if err != nil {
				fmt.Println("error retrieving channels,", err)
				return
			}

			textchannels, err := s.GuildChannels(i.GuildID)
			if err != nil {
				fmt.Println("error retrieving channels,", err)
				return
			}

			// holding first category
			var firstCategory *discordgo.Channel
			var firstVoiceCategory *discordgo.Channel

			// voice category
			for _, channel := range channels {
				if channel.Type == discordgo.ChannelTypeGuildCategory {
					// check if it has any voice channels
					for _, child := range channels {
						if child.ParentID == channel.ID && child.Type == discordgo.ChannelTypeGuildVoice {
							firstVoiceCategory = channel
							break
						}
					}
					if firstVoiceCategory != nil {
						break // stop after first category
					}
				}
			}
			// check
			if firstVoiceCategory != nil {
				fmt.Println("First Category Found:")
				fmt.Println("VoiceID:", firstVoiceCategory.ID)
			} else {
				fmt.Println("No categories found for voice")
			}

			// text category
			for _, textchannel := range textchannels {
				if textchannel.Type == discordgo.ChannelTypeGuildCategory {
					// check if it has any text channels
					for _, child := range channels {
						if child.ParentID == textchannel.ID && child.Type == discordgo.ChannelTypeGuildText {
							firstCategory = textchannel
							break
						}
					}
					if firstCategory != nil {
						break // stop after finding first category
					}
				}
			}

			// check
			if firstCategory != nil {
				fmt.Println("First Category Found:")
				fmt.Println("TextID:", firstCategory.ID)
			} else {
				fmt.Println("No categories found for text")
			}

			ParentIDChannelText = firstCategory.ID
			ParentIDVoiceText = firstVoiceCategory.ID

		}

		fmt.Printf("Selected role: %s\n", roleName)
		role := &discordgo.Role{}
		isThereCustomColor := false
		for j := 1; j < len(options); j++ {
			if options[j].Name == "custom-color" {
				isThereCustomColor = true
				break
			}
		}
		if !isThereCustomColor {
			role, err = s.GuildRoleCreate(i.GuildID, &discordgo.RoleParams{
				Name:  roleName, // Name of the new role
				Color: &color,
			})
		}

		// main cases
		for j := len(options) - 1; j > 0; j-- {

			switch options[j].Name {

			case "custom-color":
				cs := options[j].StringValue()
				if CheckHexColor(cs) {
					cs_output := cs[1:] // remove #
					value, _ := strconv.ParseInt(cs_output, 16, 32)

					color := int(value)
					fmt.Printf("Selected role: %s\n", roleName)

					role, err = s.GuildRoleCreate(i.GuildID, &discordgo.RoleParams{
						Name:  roleName,
						Color: &color,
					})

					if err != nil {
						Respond(s, i, "Failed to create role.")
						log.Println("Error creating role:", err)
						return
					}
				} else {
					Respond(s, i, "Invalid HEX color format. Example - #58a2a3 (6 hexadecimal digits).")
					return
				}

				// checkhexcolor() exists
				// if err != nil {
				//     Respond(s, i, "Value is out of range: use HEX color type (for example: #58a2a3 - 7 symbols max)")
				//     log.Println("Error creating custom color:", err)
				//     return
				// }

			case "text-category":

				ParentIDChannelText = options[j].ChannelValue(s).ID

			case "voice-category":

				ParentIDVoiceText = options[j].ChannelValue(s).ID

			default:
				fmt.Println("role dropped")
				user := options[j].UserValue(s)
				err = s.GuildMemberRoleAdd(i.GuildID, user.ID, role.ID)
				if err != nil {
					log.Println("Error adding role to member: ", err)
					return
				}

			}

		}
		Respond(s, i, "Created!") // respond for cases

		/*
			for _, userID := range len(options) {
				err := s.GuildMemberRoleAdd(i.GuildID, userID, role.ID)
				if err != nil {
					Respond(s, i, "Fuck this shits")
					return
				}
			}

			Respond(s, i, "Role created and assigned to users")
			/*
				err = s.GuildMemberRoleAdd(i.GuildID, user.ID, role.ID)
				if err != nil {
					log.Println("Error adding role to member:", err)
					return
				}
		*/

		permOverwrites := []*discordgo.PermissionOverwrite{
			{
				ID:    role.ID,
				Type:  discordgo.PermissionOverwriteTypeRole,
				Deny:  0,
				Allow: discordgo.PermissionViewChannel,
			},
			{
				ID:    i.GuildID,
				Type:  discordgo.PermissionOverwriteTypeRole,
				Deny:  discordgo.PermissionViewChannel,
				Allow: 0,
			},
		}

		channelText, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
			Name:                 roleName,
			Type:                 discordgo.ChannelTypeGuildText,
			ParentID:             ParentIDChannelText, // Tours category
			PermissionOverwrites: permOverwrites,
		})
		if err != nil {
			log.Fatalf("Cannot create channel: %v", err)
		}

		channelVoice, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
			Name:                 roleName,
			Type:                 discordgo.ChannelTypeGuildVoice,
			ParentID:             ParentIDVoiceText, // Team channels category
			PermissionOverwrites: permOverwrites,
		})
		if err != nil {
			log.Fatalf("Cannot create channel: %v", err)
		}

		fmt.Printf("Channel permissions set successfully. %s %s \n", channelText.ID, channelVoice.ID)
	} else {
		Respond(s, i, "You don't have permissions to do that!")
	}

}

func handleArchive(s *discordgo.Session, i *discordgo.InteractionCreate) {
	RespondForThinking(s, i)
	member, err := s.GuildMember(i.GuildID, i.Member.User.ID)
	if err != nil {
		fmt.Println("error retrieving member,", err)
		return
	}

	if hasRole(member.Roles, "Секретарь ЦК импрува", s, i.GuildID) {
		options := i.ApplicationCommandData().Options

		if len(options) > 0 {
			roleName := options[0].StringValue()
			fmt.Printf("Selected role: %s\n", roleName)
			err := archiveRoleMembers(s, i.GuildID, roleName, i)
			if err != nil {
				Respond(s, i, "Error: "+err.Error())
			}
			err = updateChannelPermissions(s, i.GuildID, roleName, i.ChannelID)
			if err != nil {
				Respond(s, i, "Error updating channel permissions: "+err.Error())
				return
			}
		} else {
			sendAccessibleRoles(s, i)
		}
		return
	} else {
		Respond(s, i, "WHO ARE YOU? DIDN'T READ, LOL")
	}
}

func GetUserAvatarByID(s *discordgo.Session, userID string) (string, error) {
	// fetching user
	user, err := s.User(userID)
	if err != nil {
		return "", fmt.Errorf("error fetching user: %w", err)
	}

	// building with the hash of the avatar
	avatarURL := user.AvatarURL("1024") // size 5x5 etc

	return avatarURL, nil
}

func handleRating(s *discordgo.Session, i *discordgo.InteractionCreate) {
	RespondForThinkingVisible(s, i)
	// member = _
	_, err := s.GuildMember(i.GuildID, i.Member.User.ID)
	if err != nil {
		fmt.Println("error retrieving member,", err)
		return
	}
	// fmt.Println(member.User.ID)
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		Respond(s, i, "delaem")
		return
	}
	userId := options[0].StringValue()
	fmt.Println(userId) // userId := 6560308
	url := fmt.Sprintf("http://localhost:8089/api/users/%s", userId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("error happened ", err)
		return
	}
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error happened ", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error happened ", err)
		return
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	userData, _ := data["user"].(map[string]interface{})
	userRating, ok := userData["rating"].(float64)
	ratingStr := fmt.Sprintf("%.2f", userRating)
	if ok != true {
		fmt.Println("error happened", err)
		Respond(s, i, "о нет! произошла ошибка! :(")
		return
	}
	username := userData["username"].(string)
	// data["user"]
	discord_id := userData["discord_id"].(string)
	//imageURL := "https://i.ibb.co/9sdBWTS/kth.png"

	avatarURL, err := GetUserAvatarByID(s, discord_id)
	if err != nil {
		fmt.Println("Error getting avatar:", err)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "KTH Рейтинг",
		Color: 0x00FF00,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: avatarURL,
		},

		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Details",
				Value:  fmt.Sprintf("**Username:** %s\n**Rating:** %s\n**Discord ID: **%s", username, ratingStr, discord_id),
				Inline: true,
			},
		},
	}

	editResponse := &discordgo.WebhookEdit{
		Content: nil,
		Embeds:  &[]*discordgo.MessageEmbed{embed}, // fking slice
	}

	// edit the message of the previous interaction (RespondForThinkingVisible)
	if _, err := s.InteractionResponseEdit(i.Interaction, editResponse); err != nil {
		fmt.Println("error editing response,", err)
	}
	// Respond(s, i, fmt.Sprintf("Команда еще тестируется, ваш дискорд id: %s \nKTH рейтинг игрока: %s: %f "+
	// 	"\nKTH рейтинг базируется на Skill Issue Points",
	//member.User.ID, username, userRating)
}

func archiveRoleMembers(s *discordgo.Session, guildID string, roleName string, i *discordgo.InteractionCreate) error {
	role, err := getRoleByName(s, guildID, roleName)
	if err != nil {
		return fmt.Errorf("role not found")
	}

	membersWithRole, err := getMembersWithRole(s, guildID, role.ID)
	if err != nil {
		return fmt.Errorf("couldn't find any contestants with this role")
	}

	if len(membersWithRole) == 0 {
		Respond(s, i, "There are no participants with this role")
	} else {
		memberList := "Participants with a role **" + role.Name + "**:\n"
		for _, member := range membersWithRole {
			memberList += member.User.Username + "\n"
		}
		Respond(s, i, "MemberList: "+memberList)

	}

	return nil
}

func getRoleByName(s *discordgo.Session, guildID string, roleName string) (*discordgo.Role, error) {
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		if role.Name == roleName {
			return role, nil
		}
	}
	return nil, fmt.Errorf("role not found")
}

func getMembersWithRole(s *discordgo.Session, guildID string, roleID string) ([]*discordgo.Member, error) {
	members := make([]*discordgo.Member, 0)
	guildMembers, err := s.GuildMembers(guildID, "", 1000)
	if err != nil {
		return nil, err
	}

	for _, member := range guildMembers {
		for _, r := range member.Roles {
			if r == roleID {
				members = append(members, member)
				break
			}
		}
	}

	return members, nil
}

func hasRole(roleIDs []string, roleName string, s *discordgo.Session, guildID string) bool {
	for _, roleID := range roleIDs {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			continue
		}
		if role.Name == roleName {
			return true
		}
	}
	return false
}

func sendAccessibleRoles(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channel, err := s.Channel(i.ChannelID)
	if err != nil {
		fmt.Println("error retrieving channel,", err)
		return
	}

	var accessibleRoles []string

	for _, overwrite := range channel.PermissionOverwrites {
		if overwrite.Type == discordgo.PermissionOverwriteTypeRole && overwrite.Allow&discordgo.PermissionViewChannel != 0 {
			role, err := s.State.Role(i.GuildID, overwrite.ID)
			if err != nil {
				continue
			}
			accessibleRoles = append(accessibleRoles, role.Name)
		}
	}

	if len(accessibleRoles) > 0 {
		roleList := "Roles with access to this channel:\n" + strings.Join(accessibleRoles, "\n")
		Respond(s, i, roleList)
	} else {
		Respond(s, i, "No roles have access to this channel.")
	}
}

func updateChannelPermissions(s *discordgo.Session, guildID string, roleName string, channelID string) error {
	if !strings.Contains(roleName, "[") || !strings.Contains(roleName, "]") {
		return fmt.Errorf("role name must contain '[]': %s", roleName)
	}
	role, err := getRoleByName(s, guildID, roleName)
	if err != nil {
		return err
	}

	members, err := s.GuildMembers(guildID, "", 1000)
	if err != nil {
		return fmt.Errorf("failed to retrieve guild members: %w", err)
	}

	for _, member := range members {
		for _, r := range member.Roles {
			if r == role.ID {
				err := s.ChannelPermissionSet(channelID, member.User.ID, discordgo.PermissionOverwriteTypeMember,
					discordgo.PermissionViewChannel, discordgo.PermissionSendMessages|discordgo.PermissionSendTTSMessages)
				if err != nil {
					return fmt.Errorf("failed to update permissions for user %s: %w", member.User.Username, err)
				}
			}
		}
	}
	removeVoiceChannelChoice(s, guildID, roleName, channelID)
	fmt.Printf("Updated view permissions for all members with role %s in channel %s\n", roleName, channelID)
	return nil
}

func removeVoiceChannelChoice(s *discordgo.Session, guildID string, roleName string, channelID string) {
	channels, err := s.GuildChannels(guildID)
	if err != nil {
		fmt.Println("Error retrieving channels:", err)
		return
	}

	var accessibleChannels []string

	for _, channel := range channels {
		if channel.Type == discordgo.ChannelTypeGuildVoice {
			fmt.Printf("Checking permissions for voice channel: %s\n", channel.Name)

			for _, overwrite := range channel.PermissionOverwrites {
				if overwrite.Type == discordgo.PermissionOverwriteTypeRole {
					role, err := s.State.Role(guildID, overwrite.ID)
					if err != nil {
						continue
					}
					if role.Name == roleName && overwrite.Allow&discordgo.PermissionViewChannel != 0 {
						accessibleChannels = append(accessibleChannels, channel.Name)
						fmt.Printf("Role %s has access to channel %s\n", role.Name, channel.Name)
					}
				}
			}
		}
	}

	if len(accessibleChannels) > 0 {
		fmt.Printf("Role %s has access to the following voice channels:\n", roleName)
		for _, channel := range accessibleChannels {
			fmt.Println(channel)
		}

		var options []discordgo.SelectMenuOption
		for _, channel := range accessibleChannels {
			options = append(options, discordgo.SelectMenuOption{
				Label: channel,
				Value: channel,
			})
		}
		channelsParam := strings.Join(accessibleChannels, ",")
		message := &discordgo.MessageSend{
			Content: fmt.Sprintf("Role %s has access to the following voice channels:\n%s", roleName, strings.Join(accessibleChannels, "\n")),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Delete Channels",
							CustomID: "delete_channel:" + channelsParam,
							Style:    discordgo.DangerButton,
						},
						discordgo.Button{
							Label:    "Cancel",
							CustomID: "cancel_delete_channel",
							Style:    discordgo.PrimaryButton,
						},
					},
				},
			},
		}

		_, err := s.ChannelMessageSendComplex(channelID, message)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	} else {
		fmt.Printf("Role %s does not have access to any voice channels.\n", roleName)
	}
}

func getRandomColor() int {
	rand.Seed(time.Now().UnixNano())

	r := rand.Intn(256)
	g := rand.Intn(256)
	b := rand.Intn(256)

	color := (r << 16) | (g << 8) | b
	return color
}

func handleDeletingChannels(s *discordgo.Session, i *discordgo.InteractionCreate) {
	RespondForThinking(s, i)
	member, err := s.GuildMember(i.GuildID, i.Member.User.ID)
	if err != nil {
		fmt.Println("error retrieving member,", err)
		return
	}

	if hasRole(member.Roles, "Секретарь ЦК импрува", s, i.GuildID) {

		options := i.ApplicationCommandData().Options
		var channelIDs []string

		// extracting
		for _, option := range options {
			if option.Type == discordgo.ApplicationCommandOptionChannel {
				channelIDs = append(channelIDs, option.ChannelValue(s).ID) // getting the id
			}
		}

		// deleting
		for _, channelID := range channelIDs {
			if _, err := s.ChannelDelete(channelID); err != nil {
				fmt.Println("Error deleting channel:", channelID, "Error:", err)
				continue // even if there is an error it will continue
			}
		}

		Respond(s, i, "Channels deleted successfully!")
	} else {
		Respond(s, i, "You don't have access to do that!")
	}
}
