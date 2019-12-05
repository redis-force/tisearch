package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/redis-force/tisearch/logging"
	"github.com/redis-force/tisearch/server/model"
)

type UserResponse struct {
	Results []struct {
		Gender string `json:"gender"`
		Name   struct {
			Title string `json:"title"`
			First string `json:"first"`
			Last  string `json:"last"`
		} `json:"name"`
		Location struct {
			Street struct {
				Number int    `json:"number"`
				Name   string `json:"name"`
			} `json:"street"`
			City        string `json:"city"`
			State       string `json:"state"`
			Country     string `json:"country"`
			Coordinates struct {
				Latitude  string `json:"latitude"`
				Longitude string `json:"longitude"`
			} `json:"coordinates"`
			Timezone struct {
				Offset      string `json:"offset"`
				Description string `json:"description"`
			} `json:"timezone"`
		} `json:"location"`
		Email string `json:"email"`
		Login struct {
			UUID     string `json:"uuid"`
			Username string `json:"username"`
			Password string `json:"password"`
			Salt     string `json:"salt"`
			Md5      string `json:"md5"`
			Sha1     string `json:"sha1"`
			Sha256   string `json:"sha256"`
		} `json:"login"`
		Dob struct {
			Date time.Time `json:"date"`
			Age  int       `json:"age"`
		} `json:"dob"`
		Registered struct {
			Date time.Time `json:"date"`
			Age  int       `json:"age"`
		} `json:"registered"`
		Phone string `json:"phone"`
		Cell  string `json:"cell"`
		ID    struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"id"`
		Picture struct {
			Large     string `json:"large"`
			Medium    string `json:"medium"`
			Thumbnail string `json:"thumbnail"`
		} `json:"picture"`
		Nat string `json:"nat"`
	} `json:"results"`
	Info struct {
		Seed    string `json:"seed"`
		Results int    `json:"results"`
		Page    int    `json:"page"`
		Version string `json:"version"`
	} `json:"info"`
}

var carBrands = `[{"brand": "Seat", "models": ["Alhambra", "Altea", "Altea XL", "Arosa", "Cordoba", "Cordoba Vario", "Exeo", "Ibiza", "Ibiza ST", "Exeo ST", "Leon", "Leon ST", "Inca", "Mii", "Toledo"]},
{"brand": "Renault", "models": ["Captur", "Clio", "Clio Grandtour", "Espace", "Express", "Fluence", "Grand Espace", "Grand Modus", "Grand Scenic", "Kadjar", "Kangoo", "Kangoo Express", "Koleos", "Laguna", "Laguna Grandtour", "Latitude", "Mascott", "Mégane", "Mégane CC", "Mégane Combi", "Mégane Grandtour", "Mégane Coupé", "Mégane Scénic", "Scénic", "Talisman", "Talisman Grandtour", "Thalia", "Twingo", "Wind", "Zoé"]},
{"brand": "Peugeot", "models": ["1007", "107", "106", "108", "2008", "205", "205 Cabrio", "206", "206 CC", "206 SW", "207", "207 CC", "207 SW", "306", "307", "307 CC", "307 SW", "308", "308 CC", "308 SW", "309", "4007", "4008", "405", "406", "407", "407 SW", "5008", "508", "508 SW", "605", "806", "607", "807", "Bipper", "RCZ"]},
{"brand": "Dacia", "models": ["Dokker", "Duster", "Lodgy", "Logan", "Logan MCV", "Logan Van", "Sandero", "Solenza"]},
{"brand": "Citroën", "models": ["Berlingo", "C-Crosser", "C-Elissée", "C-Zero", "C1", "C2", "C3", "C3 Picasso", "C4", "C4 Aircross", "C4 Cactus", "C4 Coupé", "C4 Grand Picasso", "C4 Sedan", "C5", "C5 Break", "C5 Tourer", "C6", "C8", "DS3", "DS4", "DS5", "Evasion", "Jumper", "Jumpy", "Saxo", "Nemo", "Xantia", "Xsara"]},
{"brand": "Opel", "models": ["Agila", "Ampera", "Antara", "Astra", "Astra cabrio", "Astra caravan", "Astra coupé", "Calibra", "Campo", "Cascada", "Corsa", "Frontera", "Insignia", "Insignia kombi", "Kadett", "Meriva", "Mokka", "Movano", "Omega", "Signum", "Vectra", "Vectra Caravan", "Vivaro", "Vivaro Kombi", "Zafira"]},
{"brand": "Alfa Romeo", "models": ["145", "146", "147", "155", "156", "156 Sportwagon", "159", "159 Sportwagon", "164", "166", "4C", "Brera", "GTV", "MiTo", "Crosswagon", "Spider", "GT", "Giulietta", "Giulia"]},
{"brand": "Škoda", "models": ["Favorit", "Felicia", "Citigo", "Fabia", "Fabia Combi", "Fabia Sedan", "Felicia Combi", "Octavia", "Octavia Combi", "Roomster", "Yeti", "Rapid", "Rapid Spaceback", "Superb", "Superb Combi"]},
{"brand": "Chevrolet", "models": ["Alero", "Aveo", "Camaro", "Captiva", "Corvette", "Cruze", "Cruze SW", "Epica", "Equinox", "Evanda", "HHR", "Kalos", "Lacetti", "Lacetti SW", "Lumina", "Malibu", "Matiz", "Monte Carlo", "Nubira", "Orlando", "Spark", "Suburban", "Tacuma", "Tahoe", "Trax"]},
{"brand": "Porsche", "models": ["911 Carrera", "911 Carrera Cabrio", "911 Targa", "911 Turbo", "924", "944", "997", "Boxster", "Cayenne", "Cayman", "Macan", "Panamera"]},
{"brand": "Honda", "models": ["Accord", "Accord Coupé", "Accord Tourer", "City", "Civic", "Civic Aerodeck", "Civic Coupé", "Civic Tourer", "Civic Type R", "CR-V", "CR-X", "CR-Z", "FR-V", "HR-V", "Insight", "Integra", "Jazz", "Legend", "Prelude"]},
{"brand": "Subaru", "models": ["BRZ", "Forester", "Impreza", "Impreza Wagon", "Justy", "Legacy", "Legacy Wagon", "Legacy Outback", "Levorg", "Outback", "SVX", "Tribeca", "Tribeca B9", "XV"]},
{"brand": "Mazda", "models": ["121", "2", "3", "323", "323 Combi", "323 Coupé", "323 F", "5", "6", "6 Combi", "626", "626 Combi", "B-Fighter", "B2500", "BT", "CX-3", "CX-5", "CX-7", "CX-9", "Demio", "MPV", "MX-3", "MX-5", "MX-6", "Premacy", "RX-7", "RX-8", "Xedox 6"]},
{"brand": "Mitsubishi", "models": ["3000 GT", "ASX", "Carisma", "Colt", "Colt CC", "Eclipse", "Fuso canter", "Galant", "Galant Combi", "Grandis", "L200", "L200 Pick up", "L200 Pick up Allrad", "L300", "Lancer", "Lancer Combi", "Lancer Evo", "Lancer Sportback", "Outlander", "Pajero", "Pajeto Pinin", "Pajero Pinin Wagon", "Pajero Sport", "Pajero Wagon", "Space Star"]},
{"brand": "Lexus", "models": ["CT", "GS", "GS 300", "GX", "IS", "IS 200", "IS 250 C", "IS-F", "LS", "LX", "NX", "RC F", "RX", "RX 300", "RX 400h", "RX 450h", "SC 430"]},
{"brand": "Toyota", "models": ["4-Runner", "Auris", "Avensis", "Avensis Combi", "Avensis Van Verso", "Aygo", "Camry", "Carina", "Celica", "Corolla", "Corolla Combi", "Corolla sedan", "Corolla Verso", "FJ Cruiser", "GT86", "Hiace", "Hiace Van", "Highlander", "Hilux", "Land Cruiser", "MR2", "Paseo", "Picnic", "Prius", "RAV4", "Sequoia", "Starlet", "Supra", "Tundra", "Urban Cruiser", "Verso", "Yaris", "Yaris Verso"]},
{"brand": "BMW", "models": ["i3", "i8", "M3", "M4", "M5", "M6", "Rad 1", "Rad 1 Cabrio", "Rad 1 Coupé", "Rad 2", "Rad 2 Active Tourer", "Rad 2 Coupé", "Rad 2 Gran Tourer", "Rad 3", "Rad 3 Cabrio", "Rad 3 Compact", "Rad 3 Coupé", "Rad 3 GT", "Rad 3 Touring", "Rad 4", "Rad 4 Cabrio", "Rad 4 Gran Coupé", "Rad 5", "Rad 5 GT", "Rad 5 Touring", "Rad 6", "Rad 6 Cabrio", "Rad 6 Coupé", "Rad 6 Gran Coupé", "Rad 7", "Rad 8 Coupé", "X1", "X3", "X4", "X5", "X6", "Z3", "Z3 Coupé", "Z3 Roadster", "Z4", "Z4 Roadster"]},
{"brand": "Volkswagen", "models": ["Amarok", "Beetle", "Bora", "Bora Variant", "Caddy", "Caddy Van", "Life", "California", "Caravelle", "CC", "Crafter", "Crafter Van", "Crafter Kombi", "CrossTouran", "Eos", "Fox", "Golf", "Golf Cabrio", "Golf Plus", "Golf Sportvan", "Golf Variant", "Jetta", "LT", "Lupo", "Multivan", "New Beetle", "New Beetle Cabrio", "Passat", "Passat Alltrack", "Passat CC", "Passat Variant", "Passat Variant Van", "Phaeton", "Polo", "Polo Van", "Polo Variant", "Scirocco", "Sharan", "T4", "T4 Caravelle", "T4 Multivan", "T5", "T5 Caravelle", "T5 Multivan", "T5 Transporter Shuttle", "Tiguan", "Touareg", "Touran"]},
{"brand": "Suzuki", "models": ["Alto", "Baleno", "Baleno kombi", "Grand Vitara", "Grand Vitara XL-7", "Ignis", "Jimny", "Kizashi", "Liana", "Samurai", "Splash", "Swift", "SX4", "SX4 Sedan", "Vitara", "Wagon R+"]},
{"brand": "Mercedes-Benz", "models": ["100 D", "115", "124", "126", "190", "190 D", "190 E", "200 - 300", "200 D", "200 E", "210 Van", "210 kombi", "310 Van", "310 kombi", "230 - 300 CE Coupé", "260 - 560 SE", "260 - 560 SEL", "500 - 600 SEC Coupé", "Trieda A", "A", "A L", "AMG GT", "Trieda B", "Trieda C", "C", "C Sportcoupé", "C T", "Citan", "CL", "CL", "CLA", "CLC", "CLK Cabrio", "CLK Coupé", "CLS", "Trieda E", "E", "E Cabrio", "E Coupé", "E T", "Trieda G", "G Cabrio", "GL", "GLA", "GLC", "GLE", "GLK", "Trieda M", "MB 100", "Trieda R", "Trieda S", "S", "S Coupé", "SL", "SLC", "SLK", "SLR", "Sprinter"]},
{"brand": "Saab", "models": ["9-3", "9-3 Cabriolet", "9-3 Coupé", "9-3 SportCombi", "9-5", "9-5 SportCombi", "900", "900 C", "900 C Turbo", "9000"]},
{"brand": "Audi", "models": ["100", "100 Avant", "80", "80 Avant", "80 Cabrio", "90", "A1", "A2", "A3", "A3 Cabriolet", "A3 Limuzina", "A3 Sportback", "A4", "A4 Allroad", "A4 Avant", "A4 Cabriolet", "A5", "A5 Cabriolet", "A5 Sportback", "A6", "A6 Allroad", "A6 Avant", "A7", "A8", "A8 Long", "Q3", "Q5", "Q7", "R8", "RS4 Cabriolet", "RS4/RS4 Avant", "RS5", "RS6 Avant", "RS7", "S3/S3 Sportback", "S4 Cabriolet", "S4/S4 Avant", "S5/S5 Cabriolet", "S6/RS6", "S7", "S8", "SQ5", "TT Coupé", "TT Roadster", "TTS"]},
{"brand": "Kia", "models": ["Avella", "Besta", "Carens", "Carnival", "Ceed", "Ceed SW", "Cerato", "K 2500", "Magentis", "Opirus", "Optima", "Picanto", "Pregio", "Pride", "Pro Ceed", "Rio", "Rio Combi", "Rio sedan", "Sephia", "Shuma", "Sorento", "Soul", "Sportage", "Venga"]},
{"brand": "Land Rover", "models": ["109", "Defender", "Discovery", "Discovery Sport", "Freelander", "Range Rover", "Range Rover Evoque", "Range Rover Sport"]},
{"brand": "Dodge", "models": ["Avenger", "Caliber", "Challenger", "Charger", "Grand Caravan", "Journey", "Magnum", "Nitro", "RAM", "Stealth", "Viper"]},
{"brand": "Chrysler", "models": ["300 C", "300 C Touring", "300 M", "Crossfire", "Grand Voyager", "LHS", "Neon", "Pacifica", "Plymouth", "PT Cruiser", "Sebring", "Sebring Convertible", "Stratus", "Stratus Cabrio", "Town & Country", "Voyager"]},
{"brand": "Ford", "models": ["Aerostar", "B-Max", "C-Max", "Cortina", "Cougar", "Edge", "Escort", "Escort Cabrio", "Escort kombi", "Explorer", "F-150", "F-250", "Fiesta", "Focus", "Focus C-Max", "Focus CC", "Focus kombi", "Fusion", "Galaxy", "Grand C-Max", "Ka", "Kuga", "Maverick", "Mondeo", "Mondeo Combi", "Mustang", "Orion", "Puma", "Ranger", "S-Max", "Sierra", "Street Ka", "Tourneo Connect", "Tourneo Custom", "Transit", "Transit", "Transit Bus", "Transit Connect LWB", "Transit Courier", "Transit Custom", "Transit kombi", "Transit Tourneo", "Transit Valnik", "Transit Van", "Transit Van 350", "Windstar"]},
{"brand": "Hummer", "models": ["H2", "H3"]},
{"brand": "Hyundai", "models": ["Accent", "Atos", "Atos Prime", "Coupé", "Elantra", "Galloper", "Genesis", "Getz", "Grandeur", "H 350", "H1", "H1 Bus", "H1 Van", "H200", "i10", "i20", "i30", "i30 CW", "i40", "i40 CW", "ix20", "ix35", "ix55", "Lantra", "Matrix", "Santa Fe", "Sonata", "Terracan", "Trajet", "Tucson", "Veloster"]},
{"brand": "Infiniti", "models": ["EX", "FX", "G", "G Coupé", "M", "Q", "QX"]},
{"brand": "Jaguar", "models": ["Daimler", "F-Pace", "F-Type", "S-Type", "Sovereign", "X-Type", "X-type Estate", "XE", "XF", "XJ", "XJ12", "XJ6", "XJ8", "XJ8", "XJR", "XK", "XK8 Convertible", "XKR", "XKR Convertible"]},
{"brand": "Jeep", "models": ["Cherokee", "Commander", "Compass", "Grand Cherokee", "Patriot", "Renegade", "Wrangler"]},
{"brand": "Nissan", "models": ["100 NX", "200 SX", "350 Z", "350 Z Roadster", "370 Z", "Almera", "Almera Tino", "Cabstar E - T", "Cabstar TL2 Valnik", "e-NV200", "GT-R", "Insterstar", "Juke", "King Cab", "Leaf", "Maxima", "Maxima QX", "Micra", "Murano", "Navara", "Note", "NP300 Pickup", "NV200", "NV400", "Pathfinder", "Patrol", "Patrol GR", "Pickup", "Pixo", "Primastar", "Primastar Combi", "Primera", "Primera Combi", "Pulsar", "Qashqai", "Serena", "Sunny", "Terrano", "Tiida", "Trade", "Vanette Cargo", "X-Trail"]},
{"brand": "Volvo", "models": ["240", "340", "360", "460", "850", "850 kombi", "C30", "C70", "C70 Cabrio", "C70 Coupé", "S40", "S60", "S70", "S80", "S90", "V40", "V50", "V60", "V70", "V90", "XC60", "XC70", "XC90"]},
{"brand": "Daewoo", "models": ["Espero", "Kalos", "Lacetti", "Lanos", "Leganza", "Lublin", "Matiz", "Nexia", "Nubira", "Nubira kombi", "Racer", "Tacuma", "Tico"]},
{"brand": "Fiat", "models": ["1100", "126", "500", "500L", "500X", "850", "Barchetta", "Brava", "Cinquecento", "Coupé", "Croma", "Doblo", "Doblo Cargo", "Doblo Cargo Combi", "Ducato", "Ducato Van", "Ducato Kombi", "Ducato Podvozok", "Florino", "Florino Combi", "Freemont", "Grande Punto", "Idea", "Linea", "Marea", "Marea Weekend", "Multipla", "Palio Weekend", "Panda", "Panda Van", "Punto", "Punto Cabriolet", "Punto Evo", "Punto Van", "Qubo", "Scudo", "Scudo Van", "Scudo Kombi", "Sedici", "Seicento", "Stilo", "Stilo Multiwagon", "Strada", "Talento", "Tipo", "Ulysse", "Uno", "X1/9"]},
{"brand": "MINI", "models": ["Cooper", "Cooper Cabrio", "Cooper Clubman", "Cooper D", "Cooper D Clubman", "Cooper S", "Cooper S Cabrio", "Cooper S Clubman", "Countryman", "Mini One", "One D"]},
{"brand": "Rover", "models": ["200", "214", "218", "25", "400", "414", "416", "620", "75"]},
{"brand": "Smart", "models": ["Cabrio", "City-Coupé", "Compact Pulse", "Forfour", "Fortwo cabrio", "Fortwo coupé", "Roadster"]}]`

var hobbies = `Aircraft Spotting
Airbrushing
Airsofting
Acting
Aeromodeling
Amateur Astronomy
Amateur Radio
Animals/pets/dogs
Archery
Arts
Astrology
Astronomy
Backgammon
Badminton
Baseball
Base Jumping
Basketball
Beach/Sun tanning
Beachcombing
Beadwork
Beatboxing
Becoming A Child Advocate
Bell Ringing
Belly Dancing
Bicycling
Bicycle Polo
Bird watching
Birding
BMX
Blacksmithing
Blogging
BoardGames
Boating
Body Building
Bonsai Tree
Bookbinding
Boomerangs
Bowling
Brewing Beer
Bridge Building
Bringing Food To The Disabled
Building A House For Habitat For Humanity
Building Dollhouses
Butterfly Watching
Button Collecting
Cake Decorating
Calligraphy
Camping
Candle Making
Canoeing
Cartooning
Car Racing
Casino Gambling
Cave Diving
Ceramics
Cheerleading
Chess
Church/church activities
Cigar Smoking
Cloud Watching
Coin Collecting
Collecting
Collecting Antiques
Collecting Artwork
Collecting Hats
Collecting Music Albums
Collecting RPM Records
Collecting Swords
Coloring
Compose Music
Computer activities
Conworlding
Cooking
Cosplay
Crafts
Crafts (unspecified)
Crochet
Crocheting
Cross-Stitch
Crossword Puzzles
Dancing
Darts
Diecast Collectibles
Digital Photography
Dodgeball
Dolls
Dominoes
Drawing
Dumpster Diving
Eating out
Educational Courses
Electronics
Embroidery
Entertaining
Exercise
Falconry
Fast cars
Felting
Fencing
Fire Poi
Fishing
Floorball
Floral Arrangements
Fly Tying
Football
Four Wheeling
Freshwater Aquariums
Frisbee Golf – Frolf
Games
Gardening
Garage Saleing
Genealogy
Geocaching
Ghost Hunting
Glowsticking
Gnoming
Going to movies
Golf
Go Kart Racing
Grip Strength
Guitar
Gunsmithing
Gun Collecting
Gymnastics
Gyotaku
Handwriting Analysis
Hang gliding
Herping
Hiking
Home Brewing
Home Repair
Home Theater
Horse riding
Hot air ballooning
Hula Hooping
Hunting
Iceskating
Illusion
Impersonations
Internet
Inventing
Jet Engines
Jewelry Making
Jigsaw Puzzles
Juggling
Keep A Journal
Jump Roping
Kayaking
Kitchen Chemistry
Kites
Kite Boarding
Knitting
Knotting
Lasers
Lawn Darts
Learn to Play Poker
Learning A Foreign Language
Learning An Instrument
Learning To Pilot A Plane
Leathercrafting
Legos
Letterboxing
Listening to music
Locksport
Lacrosse
Macramé
Magic
Making Model Cars
Marksmanship
Martial Arts
Matchstick Modeling
Meditation
Microscopy
Metal Detecting
Model Railroading
Model Rockets
Modeling Ships
Models
Motorcycles
Mountain Biking
Mountain Climbing
Musical Instruments
Nail Art
Needlepoint
Owning An Antique Car
Origami
Painting
Paintball
Papermaking
Papermache
Parachuting
Paragliding or Power Paragliding
Parkour
People Watching
Photography
Piano
Pinochle
Pipe Smoking
Planking
Playing music
Playing team sports
Pole Dancing
Pottery
Powerboking
Protesting
Puppetry
Pyrotechnics
Quilting
Racing Pigeons
Rafting
Railfans
Rapping
R/C Boats
R/C Cars
R/C Helicopters
R/C Planes
Reading
Reading To The Elderly
Relaxing
Renaissance Faire
Renting movies
Rescuing Abused Or Abandoned Animals
Robotics
Rock Balancing
Rock Collecting
Rockets
Rocking AIDS Babies
Roleplaying
Running
Saltwater Aquariums
Sand Castles
Scrapbooking
Scuba Diving
Self Defense
Sewing
Shark Fishing
Skeet Shooting
Skiing
Shopping
Singing In Choir
Skateboarding
Sketching
Sky Diving
Slack Lining
Sleeping
Slingshots
Slot Car Racing
Snorkeling
Snowboarding
Soap Making
Soccer
Socializing with friends/neighbors
Speed Cubing
Spelunkering
Spending time with family_kids
Stamp Collecting
Storm Chasing
Storytelling
String Figures
Surfing
Surf Fishing
Survival
Swimming
Tatting
Taxidermy
Tea Tasting
Tennis
Tesla Coils
Tetris
Texting
Textiles
Tombstone Rubbing
Tool Collecting
Toy Collecting
Train Collecting
Train Spotting
Traveling
Treasure Hunting
Trekkie
Tutoring Children
TV watching
Ultimate Frisbee
Urban Exploration
Video Games
Violin
Volunteer
Walking
Warhammer
Watching sporting events
Weather Watcher
Weightlifting
Windsurfing
Wine Making
Wingsuit Flying
Woodworking
Working In A Food Pantry
Working on cars
World Record Breaking
Wrestling
Writing
Writing Music
Writing Songs
Yoga
YoYo
Ziplining
Zumba`

func main() {
	brands := make([]struct {
		Brand string `json:brand`
	}, 0)
	json.Unmarshal([]byte(carBrands), &brands)
	hobbiesCorpus := strings.Split(hobbies, "\n")
	var words []string
	for _, b := range brands {
		words = append(words, strings.ToLower(b.Brand))
	}
	for _, h := range hobbiesCorpus {
		words = append(words, strings.ToLower(strings.Replace(strings.TrimSpace(h), " ", "_", -1)))
	}
	fmt.Println(strings.Join(words, "\",\""))
	choiceCarBrand := func() string {
		n := rand.Intn(7)
		occur := make(map[string]struct{})
		var cars []string
		for i := 0; i < n; i++ {
			j := rand.Intn(len(brands))
			car := strings.ToLower(brands[j].Brand)
			if _, ok := occur[car]; !ok {
				cars = append(cars, "car_"+car)
			}
		}
		return strings.Join(cars, " ")
	}
	choiceHobbies := func() string {
		n := rand.Intn(10)
		occur := make(map[string]struct{})
		var hs []string
		for i := 0; i < n; i++ {
			j := rand.Intn(len(hobbiesCorpus))
			hobby := strings.ToLower(strings.Replace(strings.TrimSpace(hobbiesCorpus[j]), " ", "_", -1))
			if _, ok := occur[hobby]; !ok {
				hs = append(hs, "interest_"+hobby)
			}
		}
		return strings.Join(hs, " ")
	}

	resp, err := http.Get("https://randomuser.me/api/?results=1000")
	if err != nil {
		logging.Errorf("get users error %s", err)
		panic(err)
	}
	var userResponse UserResponse
	err = json.NewDecoder(resp.Body).Decode(&userResponse)
	if err != nil {
		panic(err)
	}
	dbDSN := "root:@tcp(10.9.118.254:4000)/tisearch?charset=utf8&timeout=1s&parseTime=true"
	db, err := gorm.Open("mysql", dbDSN)
	db = db.Debug()
	if err != nil {
		panic(err)
	}
	id := int64(1)
	for _, u := range userResponse.Results {
		var searchUser model.User
		searchUser.Name = u.Name.Title + " " + u.Name.First + " " + u.Name.Last
		searchUser.Location = u.Location.Country + " " + u.Location.City
		searchUser.Picture = u.Picture.Large
		searchUser.Labels = choiceHobbies() + " " + choiceCarBrand()
		searchUser.Birthday = u.Dob.Date
		searchUser.Coordinates = u.Location.Coordinates.Latitude + "," + u.Location.Coordinates.Longitude
		searchUser.Gender = 1
		if u.Gender == "femal" {
			searchUser.Gender = 2
		}
		id++
		if err := db.Create(&searchUser).Error; err != nil {
			logging.Errorf("save user error %s", err)
			panic(err)
		}
	}

}
