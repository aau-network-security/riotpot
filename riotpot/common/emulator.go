/* 
Package environ provides structures for the emulation of services
*/
package emulator
import (
	"plugin"
	"log"
)

type Emulator struct {
	ID int
	Name string
}

type Emulators struct{
	ID int
	emulators []emulator.Emulator
}

func (es *Emulators) register(emulator Emulator){
	"""Method used to append a new emulator to the list of this object"""
	es.emulators = append(es.emulators, emulator)
}

func register_emulators(validated_emulators) *Emulators{
	log.Print("Registering Emulators...")

	// based on: https://echorand.me/posts/getting-started-with-golang-plugins/
	for _, emu := range validated_emulators{
		// define the path of where the plugin is installed
		var path string = "/emulators/" + emu + "/plugin.so"

		// open the file
		p, err := plugin.Open(path)
		if err != nil {
			log.Fatal(err)
		}

		rf, err := p.Lookup("Register")
		if err != nil {
			panic(err)
		}
		// run the register function 
		rf()
	}

	log.Print("All emulators registered successfully")
	return &Emulators{emulators: emulators}
}