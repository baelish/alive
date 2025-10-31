package server

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/baelish/alive/api"
)

var animals = [...]string{
	"Aardvark",
	"African Elephant",
	"African Tree Pangolin",
	"Albatross",
	"Alligator",
	"Alpaca",
	"Anaconda",
	"Angel Fish",
	"Ant",
	"Anteater",
	"Antelope",
	"Arab horse",
	"Archer Fish",
	"Armadillo",
	"Asian Elephant",
	"Atlantic Puffin",
	"Aye-Aye",
	"Baboon",
	"Badger",
	"Bald Eagle",
	"Bandicoot",
	"Bangle Tiger",
	"Barnacle",
	"Barracuda",
	"Basilisk",
	"Bass",
	"Basset Hound",
	"Bat",
	"Bearded Dragon",
	"Beaver",
	"Bee",
	"Beetle",
	"Beluga Whale",
	"Big-horned sheep",
	"Billy goat",
	"Bird of paradise",
	"Bird",
	"Bison",
	"Black Bear",
	"Black Fly",
	"Black Footed Rhino",
	"Black Rhino",
	"Black Widow Spider",
	"Blackbird",
	"Blowfish",
	"Blue Jay",
	"Blue Whale",
	"Boa",
	"Boar",
	"Bob-Cat",
	"Bonobo",
	"Border Collie",
	"Bornean Orang-utan",
	"Bottle-Nose dolphin",
	"Boxer dog",
	"Brown Bear",
	"Buck",
	"Budgie",
	"Buffalo",
	"Bull Mastiff",
	"Bull frog",
	"Bull",
	"Butterfly",
	"Buzzard",
	"Caiman lizard",
	"Camel",
	"Canary",
	"Caribou",
	"Carp",
	"Cat",
	"Caterpillar",
	"Catfish",
	"Cattle",
	"Centipede",
	"Chameleon",
	"Cheetah",
	"Chicken",
	"Chihuahua",
	"Chimpanzee",
	"Chinchilla",
	"Chipmunk",
	"Chupacabra",
	"Clam",
	"Clown Fish",
	"Cobra",
	"Cockatiel",
	"Cockatoo",
	"Cocker Spaniel",
	"Cockroach",
	"Cod",
	"Coho",
	"Common Dolphin",
	"Common seal",
	"Corn Snake",
	"Cougar",
	"Cow",
	"Coyote",
	"Crab",
	"Crane",
	"Crawfish",
	"Cray fish",
	"Cricket",
	"Crocodile",
	"Crow",
	"Cuckoo bird",
	"Cuttle fish",
	"Dacshund",
	"Dalmation",
	"Damsel fly",
	"Dart Frog",
	"Deer",
	"Devi Fish (Giant Sting ray)",
	"Diamond back rattler",
	"Dik-dik",
	"Dingo",
	"Dinosaur",
	"Doberman Pinscher",
	"Dodo bird",
	"Dog",
	"Dolly Varden",
	"Dolphin",
	"Donkey",
	"Door mouse",
	"Dormouse",
	"Dove",
	"Draft horse",
	"Dragonfly",
	"Drake",
	"Du-gong",
	"Duck",
	"Duckbill Platypus",
	"Dung beetle",
	"Eagle",
	"Earthworm",
	"Earwig",
	"Echidna",
	"Eclectus",
	"Eel",
	"Egret",
	"Elephant Seal",
	"Elephant",
	"Elk",
	"Emu",
	"Erne",
	"Eurasian Lynx",
	"Falcon",
	"Ferret",
	"Finch",
	"Firefly",
	"Fish",
	"Flamingo",
	"Flatworm",
	"Fly",
	"Fox",
	"Frog",
	"Gazelle",
	"Giant Anteater",
	"Giant panda",
	"Giraffe",
	"Gnat",
	"Goat",
	"Goose",
	"Gopher",
	"Gorilla",
	"Grasshopper",
	"Great White Shark",
	"Green fly",
	"Grey Whale",
	"Groundhog",
	"Hammerhead shark",
	"Hare",
	"Hawk",
	"Hedgehog",
	"Heron",
	"Herring",
	"Hippopotamus",
	"Horse",
	"Hyena",
	"Hyrax",
	"Iguana",
	"Iguanodon",
	"Impala",
	"Inchworm",
	"Insect",
	"Jackal",
	"Jackrabbit",
	"Jaguar",
	"Jellyfish",
	"June bug",
	"Kangaroo",
	"Killer Whale",
	"King Cobra",
	"Kingfisher",
	"Koala",
	"Komodo Dragon",
	"Kookaburra",
	"Krill",
	"Lama",
	"Lamb",
	"Lancelet",
	"Leech",
	"Lemming",
	"Lemur",
	"Leopard",
	"Lice",
	"Lion",
	"Lionfish",
	"Llama",
	"Lobster",
	"Lynx",
	"Man-Of-War",
	"Manatee",
	"Mantis",
	"Marmot",
	"Marsupials",
	"Meerkat",
	"Mink",
	"Mole",
	"Mollusks",
	"Monarch Butterfly",
	"Mongoose",
	"Monkey",
	"Moose",
	"Mountain Lion",
	"Mouse",
	"Mule",
	"Muskox",
	"Muskrat",
	"Naked Mole Rat",
	"Narwhal",
	"Nautilus",
	"Newt",
	"Ocelot",
	"Octopus",
	"Opossum",
	"Orangutan",
	"Orca",
	"Osprey",
	"Ostrich",
	"Otter",
	"Owl",
	"Ox",
	"Panda",
	"Panther",
	"Peacock",
	"Pelican",
	"Penguin",
	"Pig",
	"Pigeon",
	"Platypus",
	"Polar Bear",
	"Porcupine",
	"Prawn",
	"Praying Mantis",
	"Puma",
	"Quail",
	"Quetzal",
	"Rabbit",
	"Raccoon",
	"Rat",
	"Ray",
	"Reindeer",
	"Rhino",
	"Rhinoceros",
	"Ringworm",
	"Robin",
	"Rooster",
	"Roundworm",
	"Salmon",
	"Sandpiper",
	"Scallop",
	"Scorpion",
	"Sea Lion",
	"Sea anemone",
	"Sea urchin",
	"Seahorse",
	"Seal",
	"Shark",
	"Sheep",
	"Shrimp",
	"Siberian Husky",
	"Siberian Tiger",
	"Skunks",
	"Slender Loris",
	"Sloth bear",
	"Sloth",
	"Slugs",
	"Snails",
	"Snake",
	"Snow Fox",
	"Snow Hare",
	"Snow Leopard",
	"Somali Wild Ass",
	"Spectacled Bear",
	"Sponge",
	"Squid",
	"Squirrel",
	"Starfish",
	"Stork",
	"Swan",
	"Swordfish",
	"Tadpole",
	"Tamarin",
	"Tapeworm",
	"Tapir",
	"Tarantula",
	"Tarpan",
	"Tasmanian Devil",
	"Tazmanian devil",
	"Tazmanian tiger",
	"Terrapin",
	"Tick",
	"Tiger shark",
	"Tiger",
	"Tortoise",
	"Trout",
	"Tuna",
	"Turkey",
	"Turtle",
	"Uakari",
	"Umbrella bird",
	"Urchin",
	"Urutu",
	"Vampire bat",
	"Velociraptor",
	"Velvet worm",
	"Vervet",
	"Vicuna",
	"Viper Fish",
	"Viper",
	"Vole",
	"Vulture",
	"Wallaby",
	"Walrus",
	"Warbler",
	"Warthog",
	"Wasp",
	"Water Buffalo",
	"Water Dragons",
	"Weasel",
	"Weevil",
	"Whale Shark",
	"Whale",
	"Whippet",
	"White Rhino",
	"White tailed dear",
	"Whooper",
	"Whooping Crane",
	"Widow Spider",
	"Wildcat",
	"Wildebeest",
	"Wolf Spider",
	"Wolf",
	"Wolverine",
	"Wombat",
	"Woodchuck",
	"Woodpecker",
	"Wren",
	"X-ray fish",
	"Yak",
	"Yeti",
	"Yorkshire terrier",
	"Zander",
	"Zebra Dove",
	"Zebra Finch",
	"Zebra",
	"Zebu",
	"Zorilla",
}

func createRandomBox() {
	var newBox api.Box
	newBox.Name = animals[rand.Intn(len(animals))]
	newBox.Size = api.BoxSize(rand.Intn(int(api.Xlarge)-int(api.Dot)+1) + int(api.Dot))
	info := map[string]string{}
	info["foo"] = "bar"
	info["boo"] = "hoo"
	newBox.Info = &info
	newBox.Status = api.Grey
	addBox(newBox)
}

func runDemo(ctx context.Context) {
	if options.Debug {
		logger.Info("Starting demo routine")
	}

	const (
		minDemoBoxes = 10
		maxDemoBoxes = 60
		pause        = time.Duration(50 * time.Millisecond)
	)
	// Generate a constant stream of events that get pushed
	// into the Broker's messages channel and are then broadcast
	// out to any clients that are attached.
	x := 0
	var event api.Event

	// Create a box if there are none.
	if boxStore.Len() == 0 {
		createRandomBox()
	}

	for {
		boxCount := boxStore.Len()
		// Get all boxes once per iteration to avoid repeated allocations
		allBoxes := boxStore.GetAll()
		max := boxCount - 1
		if max < 1 {
			max = 1
		}
		switch e := rand.Intn(100); {
		case e < 5: // Create a box
			if boxCount < maxDemoBoxes {
				createRandomBox()
			}
		case e < 10: // Delete a box
			if boxCount > minDemoBoxes {
				// Get a random box to delete
				if len(allBoxes) > 0 {
					deleteBox(allBoxes[rand.Intn(len(allBoxes))].ID, true)
				}
			}
		case e < 20: // Update a box with a random event
			y := rand.Intn(max)

			switch rand.Intn(3) {
			case 0:
				event.Status = api.Red
				event.Message = "PANIC! Red Alert"
			case 1:
				event.Status = api.Amber
				event.Message = "OH NOES! Something's not quite right"
			case 2:
				event.Status = api.Grey
				event.Message = "Meh not sure what to do now...."
			}

			if y < len(allBoxes) {
				event.ID = allBoxes[y].ID
				update(event)
			}
		case e < 25: // Set Max TBU to small number
			event.MaxTBU.Duration = time.Second * 4
			event.ExpireAfter.Duration = 0
			event.Message = "Adding 4s MaxTBU"
			event.Status = api.Green
			if len(allBoxes) > 0 {
				event.ID = allBoxes[rand.Intn(max)].ID
				update(event)
			}

		case e < 30: // Set Max TBU to small number
			event.MaxTBU.Duration = 0
			event.ExpireAfter.Duration = 5 * time.Second
			event.Message = "Expiring box in 5s"
			event.Status = api.Grey
			if len(allBoxes) > 0 {
				event.ID = allBoxes[rand.Intn(max)].ID
				update(event)
			}
		default:
			x++
			if x >= len(allBoxes) {
				x = 0
			}

			if len(allBoxes) > 0 && x < len(allBoxes) {
				id := allBoxes[x].ID
				// Create a little message to send to clients,
				// including the current time.
				t := time.Now()
				ft := t.Format(timeFormat)

				event.ID = id
				event.Status = api.Green
				event.Message = fmt.Sprintf("the time is %s", ft)
				update(event)
			}

		}

		select {
		case <-ctx.Done():
			return

		case <-time.After(pause):
		}
	}
}
