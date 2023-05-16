package customize

import (
	"fmt"

	"github.com/robertcharca/skittyc/kittyc/kfeatures"
	"github.com/robertcharca/skittyc/prompts"
	"github.com/spf13/cobra"
)

var (
	//Local flags
	setting bool	
	changing bool
)

var fontCmd = &cobra.Command{
	Use: "fonts",
	Short: "'fonts' subcommand",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("'fonts' working")
		
		if changing == true {
			prompt, res := prompts.HandleFontChangeValues()
			fmt.Printf("Prompt message: %s\n Prompt result: %s", prompt, res)
			kfeatures.ChangingFontValues(prompt, res)
		}

		if setting == true {
			fmt.Println("Set")
		}			
	},
}

func init() {	
	//Adding the command `fonts` to `customize`.
	CustomizeCmd.AddCommand(fontCmd)
	
	//Flags
	/*
		An argument is called using the flag variable and setting your command variable
		with the "Flags()" method. 
	*/
	// Flag for changing font values (size, bold, italic)
	fontCmd.Flags().BoolVarP(&changing, "change", "c", false, "Changing font values.")
	// Flag for setting a new font 
	fontCmd.Flags().BoolVarP(&setting, "set", "s", false, "Setting a new font.")	
}
