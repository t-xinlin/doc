package main

import "fmt"

type USB interface {
	Name() string
	Connector
}

type Connector interface {
	Connect()
}

type PhoneConnector struct {
	name string
}

func (pc *PhoneConnector) Name() string {
	return pc.name
}

func (pc *PhoneConnector) Connect() {
	pc.name = "New Name"
	fmt.Printf("Connect to %s ok \n", pc.Name())
}

func DisConnect(usb USB) {
	if pc, ok := usb.(*PhoneConnector); ok {
		fmt.Println("DisConnect from ", pc.Name())
		return
	}
	fmt.Println("Unkonw device")
}

func main() {
	var a USB
	a = &PhoneConnector{"Iphone"}
	a.Connect()
	DisConnect(a)

	fmt.Println("-----------------")

	var phone Phone
	phone = &Nokia{Name: "Nokia"}
	phone.Call()
	phone = &Iphone{Name: "Iphone"}
	phone.Call()

}

type Phone interface {
	Call()
}

type Nokia struct {
	Name string
}

func (nokia *Nokia) Call() {
	fmt.Println("I am ", nokia.Name)
}

type Iphone struct {
	Name string
}

func (iphone *Iphone) Call() {
	fmt.Println("I am ", iphone.Name)
}
