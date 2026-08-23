package main

import "fmt"

type Engine struct {
	Model      string
	HorsePower int
}

func (e *Engine) Start() {
	fmt.Println("Engine started!")
}

type GPS struct {
	Model string
}

type Car struct {
	Model  string
	Engine // Composition: Car has an Engine
	GPS    // Composition: Car has a GPS
}

func (c *Car) Drive() {
	fmt.Printf("Driving the %s\n", c.Model)
}

func main() {
	// Go is OOP, but does not have concept of inheritance
	// Instead, we can use composition to achieve similar results

	myCar := Car{
		Model: "Toyota Camry",
		Engine: Engine{
			HorsePower: 200,
			Model:      "V6 Engine",
		},
		GPS: GPS{
			Model: "Garmin",
		},
	}

	// fmt.Println("Car Model:", myCar.Model)
	// fmt.Println("Engine HorsePower:", myCar.HorsePower) // Accessing Engine's HorsePower directly from Car
	// myCar.Start()
	// myCar.Drive()
	// fmt.Println(myCar.Model, myCar.Engine.Model)

	fmt.Println("Car Model:", myCar.Model)
	fmt.Println("Engine HorsePower:", myCar.HorsePower)
	fmt.Println("GPS Model:", myCar.GPS.Model)
}
