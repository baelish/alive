package server

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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

func ptr[T any](v T) *T { return &v }

// demoLayout is the set of boxes the demo starts with.
var demoLayout = []struct {
	size  api.BoxSize
	count int
}{
	{api.Xlarge, 1},
	{api.Dlarge, 1},
	{api.Large, 1},
	{api.Dmedium, 1},
	{api.Medium, 1},
	{api.Dsmall, 3},
	{api.Small, 5},
	{api.Dmicro, 2},
	{api.Micro, 8},
	{api.Dot, 32},
}

// seedDemoBoxes creates the boxes in demoLayout. A size with a single box
// always gets a stat, where there are several about half of them do.
func seedDemoBoxes() {
	for _, l := range demoLayout {
		withStat := 1
		if l.count > 1 {
			withStat = l.count / 2
		}
		for i := 0; i < l.count; i++ {
			createDemoBox(l.size, i < withStat)
		}
	}
}

// randomStat returns a short value in one of a few formats to show off the box
// stat field.
func randomStat() string {
	switch rand.Intn(4) {
	case 0:
		total := rand.Intn(20) + 1
		return fmt.Sprintf("%d/%d", rand.Intn(total+1), total)
	case 1:
		return fmt.Sprintf("%d%%", rand.Intn(101))
	case 2:
		return fmt.Sprintf("%dms", rand.Intn(500)+1)
	default:
		return fmt.Sprintf("%.1fk", rand.Float64()*100)
	}
}

// nextStat returns the next value for a demo stat, keeping the format of the
// current one (a "3/10" stays out of 10, a "120ms" stays in ms). Half the time
// the value is left as it is, otherwise it drifts a little from where it was.
func nextStat(current string) string {
	if rand.Intn(2) == 0 {
		return current
	}

	var n, total int
	if _, err := fmt.Sscanf(current, "%d/%d", &n, &total); err == nil {
		return fmt.Sprintf("%d/%d", clamp(n+rand.Intn(3)-1, 0, total), total)
	}
	if v, ok := strings.CutSuffix(current, "ms"); ok {
		if n, err := strconv.Atoi(v); err == nil {
			return fmt.Sprintf("%dms", clamp(n+rand.Intn(41)-20, 1, 999))
		}
	}
	if v, ok := strings.CutSuffix(current, "%"); ok {
		if n, err := strconv.Atoi(v); err == nil {
			return fmt.Sprintf("%d%%", clamp(n+rand.Intn(11)-5, 0, 100))
		}
	}
	if v, ok := strings.CutSuffix(current, "k"); ok {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return fmt.Sprintf("%.1fk", max(f*(0.95+rand.Float64()*0.1), 0))
		}
	}

	// Not a format the demo knows, start again with a new stat.
	return randomStat()
}

func clamp(v, lo, hi int) int {
	return min(max(v, lo), hi)
}

func createDemoBox(size api.BoxSize, withStat bool) {
	var newBox api.Box
	newBox.Name = animals[rand.Intn(len(animals))]
	newBox.Size = size
	info := map[string]string{}
	info["foo"] = "bar"
	info["boo"] = "hoo"
	newBox.Info = &info
	newBox.Status = api.Grey
	if withStat {
		newBox.Stat = randomStat()
	}
	addBox(newBox)
}

// createMissingDemoBox replaces a box of a size that has dropped below its
// count in demoLayout, so the demo keeps its shape as boxes are deleted.
func createMissingDemoBox(boxes []api.Box) {
	counts := map[api.BoxSize]int{}
	for _, b := range boxes {
		counts[b.Size]++
	}

	var missing []int
	for i, l := range demoLayout {
		if counts[l.size] < l.count {
			missing = append(missing, i)
		}
	}
	if len(missing) == 0 {
		return
	}

	l := demoLayout[missing[rand.Intn(len(missing))]]
	createDemoBox(l.size, l.count == 1 || rand.Intn(2) == 0)
}

func runDemo(ctx context.Context) {
	if options.Debug {
		logger.Info("Starting demo routine")
	}

	const (
		minDemoBoxes = 10
		maxDemoBoxes = 60
		pause        = 300 * time.Millisecond
	)
	// Generate a constant stream of events that get pushed
	// into the Broker's messages channel and are then broadcast
	// out to any clients that are attached.
	x := 0

	// Start with a fixed set of boxes if there are none.
	if boxStore.Len() == 0 {
		seedDemoBoxes()
	}

	for {
		boxCount := boxStore.Len()
		// Get all boxes once per iteration to avoid repeated allocations
		allBoxes := boxStore.GetAll()
		max := boxCount - 1
		if max < 1 {
			max = 1
		}
		// Start from an empty event each time so settings such as ExpireAfter
		// only apply to the box they were meant for.
		var event api.Event
		switch e := rand.Intn(100); {
		case e < 15: // Replace a box missing from the demo layout
			if boxCount < maxDemoBoxes {
				createMissingDemoBox(allBoxes)
			}
		case e < 20: // Delete a box
			if boxCount > minDemoBoxes {
				// Get a random box to delete
				if len(allBoxes) > 0 {
					deleteBox(allBoxes[rand.Intn(len(allBoxes))].ID, true)
				}
			}
		case e < 30: // Update a box with a random event
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
		case e < 35: // Set Max TBU to small number and unset expiry
			event.MaxTBU = ptr(api.Duration(pause))
			event.ExpireAfter = ptr(api.Duration(0))
			event.Message = fmt.Sprintf("Setting MaxTBU to %s", pause)
			event.Status = api.Green
			if len(allBoxes) > 0 {
				event.ID = allBoxes[rand.Intn(max)].ID
				update(event)
			}

		case e < 40: // Unset Max TBU and set expiry time
			event.MaxTBU = ptr(api.Duration(0))
			event.ExpireAfter = ptr(api.Duration(5 * pause))
			event.Message = fmt.Sprintf("Expiring box in %s ", 5*pause)
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
				if allBoxes[x].Stat != "" {
					event.Stat = ptr(nextStat(allBoxes[x].Stat))
				}
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
