package main

//wdym identical
import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

type Trains []Train

type Train struct {
	TrainID            int
	DepartureStationID int
	ArrivalStationID   int
	Price              float32
	ArrivalTime        time.Time
	DepartureTime      time.Time
}

func (t *Train) UnmarshalJSON(data []byte) error {

	var objects map[string]*json.RawMessage //we put all json data into a map, keys represent the names of variables and *json.RawMessage is raw json data

	err := json.Unmarshal(data, &objects) //unmarshalling
	if err != nil {
		return err
	}
	//then we normally unmarshal the data into all the Train struct variables except for ArrivalTime and Departure time.
	err = json.Unmarshal(*objects["trainId"], &t.TrainID)
	if err != nil {
		return err
	}
	err = json.Unmarshal(*objects["departureStationId"], &t.DepartureStationID)
	if err != nil {
		return err
	}
	err = json.Unmarshal(*objects["arrivalStationId"], &t.ArrivalStationID)
	if err != nil {
		return err
	}
	err = json.Unmarshal(*objects["price"], &t.Price)
	if err != nil {
		return err
	}
	var Arr, Dep string //we use string variables to get the lines from the json file
	err = json.Unmarshal(*objects["arrivalTime"], &Arr)
	if err != nil {
		return err
	}
	err = json.Unmarshal(*objects["departureTime"], &Dep)
	if err != nil {
		return err
	}
	var h, m, s int
	_, err = fmt.Sscanf(Arr, "%d:%d:%d", &h, &m, &s) //then we convert them into the varaibles h(ours), m(inutes),and s(econds)
	if err != nil {
		return err
	}
	t.ArrivalTime = time.Date(0, time.January, 1, //and then we set the Arrival and Departure time!
		h, m, s, 0,
		time.UTC)
	_, err = fmt.Sscanf(Dep, "%d:%d:%d", &h, &m, &s)
	if err != nil {
		return err
	}
	t.DepartureTime = time.Date(0, time.January, 1,
		h, m, s, 0,
		time.UTC)
	return nil
}
func main() {

	byteValue, err := ioutil.ReadFile("data.json") //Opening file
	if err != nil {
		log.Fatal("Error when opening file:", err)
	}

	var trains Trains

	err = json.Unmarshal(byteValue, &trains) //Unmarshalling the file into the variable trains
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	stationDepNums := make(map[int][]int) //maps we'll use to find trains that match our departure and arrival stations
	stationArrNums := make(map[int][]int)

	for i := 0; i < len(trains); i++ {
		stationDepNums[trains[i].DepartureStationID] = append(stationDepNums[trains[i].DepartureStationID], i) //adds trains to matching maps
		stationArrNums[trains[i].ArrivalStationID] = append(stationArrNums[trains[i].ArrivalStationID], i)     //i.e if DepartureStationID = 10 in a Train varaible, the location of the train within Trains gets added to a slice in stationDepNums corresponding to the key "10"
	}

	var departureStation string
	var arrivalStation string
	var criteria string

	fmt.Println("Enter departure station:") //departure
	fmt.Scanln(&departureStation)

	if departureStation == "" {
		log.Fatal(errors.New("empty departure station"))
	}
	if _, err := strconv.Atoi(departureStation); err != nil {
		log.Fatal(errors.New("bad departure station input"))
	}

	fmt.Println("Enter arrival station:") //arrival
	fmt.Scanln(&arrivalStation)           //Yes, i know that the varaible named "departureStation" is used to get the arrival station ID. this is because i noticed that

	if arrivalStation == "" {
		log.Fatal(errors.New("empty arrival station"))
	}
	if _, err := strconv.Atoi(arrivalStation); err != nil {
		log.Fatal(errors.New("bad arrival station input"))
	}

	fmt.Println("Enter criteria for search(price, arrival-time, departure-time):")
	fmt.Scanln(&criteria)

	result, err := FindTrains(departureStation, arrivalStation, criteria, trains, stationArrNums, stationDepNums)
	if err != nil {
		log.Fatal(err)
	}
	for _, val := range result {
		fmt.Printf(
			"{TrainID: %v, DepartureStationID: %v, ArrivalStationID: %v, Price: %v, ArrivalTime: %v, DepartureTime: %v}\n",
			val.TrainID, val.DepartureStationID, val.ArrivalStationID, val.Price, val.ArrivalTime, val.DepartureTime)
	}
}

func FindTrains(d string, a string, criteria string, trains Trains, Arrs map[int][]int, Deps map[int][]int) (Trains, error) {
	departureStation, _ := strconv.Atoi(d)
	arrivalStation, _ := strconv.Atoi(a)
	ourDeps := Deps[departureStation]
	ourArrs := Arrs[arrivalStation] //we get locations of trains that match our departureStation and arrivalStation

	var needed Trains //trains we *need* to sort
	for _, val := range ourArrs {
		for i, val2 := range ourDeps {
			if val2 == val {
				needed = append(needed, trains[ourDeps[i]]) //when a train from ourDeps is encounterred within ourArrs, it gets added to the variable needed
			}
		}
	}
	if len(needed) == 0 { //if we have no matching trains, we return nothing
		return nil, nil
	}
	var sorted Trains
	switch criteria {
	case "price": // sort by price, lowest to highest
		var mintrain Train
		var minval int
		for i := 0; i < len(needed); i++ {

			minval = 0
			mintrain.Price = 0.00 //we only need to compare by price, so the only struct variable within mintrain that we'll be directly chaning is Price. Same with others

			for i2, val := range needed {
				if val.Price > mintrain.Price {
					mintrain = val
					minval = i2
				}
			}

			needed[minval] = needed[len(needed)-1]
			needed = needed[:len(needed)-1]
			sorted = append(Trains{mintrain}, sorted...)

		}
		return sorted, nil

	case "arrival-time": //sort by arrival time, lowest to highest
		var mintrain Train
		var minval int
		for i := 0; i < len(needed); i++ {

			minval = 0
			mintrain.ArrivalTime = time.Date(0, 01, 01, 99, 99, 99, 0, time.UTC)

			for i2, val := range needed {
				if val.ArrivalTime.Before(mintrain.ArrivalTime) {
					mintrain = val
					minval = i2
				}
			}

			needed[minval] = needed[len(needed)-1]
			needed = needed[:len(needed)-1]
			sorted = append(sorted, mintrain)

		}
		return sorted, nil

	case "departure-time": //sort by departure time, lowest to highest

		var mintrain Train
		var minval int

		for i := 0; i < len(needed); i++ {
			minval = 0
			mintrain.DepartureTime = time.Date(0, 01, 01, 99, 99, 99, 0, time.UTC)
			for i2, val := range needed {
				if val.DepartureTime.Before(mintrain.DepartureTime) {
					mintrain = val
					minval = i2
				}
			}

			needed[minval] = needed[len(needed)-1]
			needed = needed[:len(needed)-1]
			sorted = append(sorted, mintrain)

		}
		return sorted, nil

	default:
		return nil, errors.New("unsupported criteria")

	}
}
