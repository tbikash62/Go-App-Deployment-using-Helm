package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
)

// Mass represents the density of the material
type Mass struct {
	Density float64
}

// MassVolume interface that both Sphere and Cube will implement
type MassVolume interface {
	density() float64
	volume(dimension float64) float64
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////
// Sphere represents a sphere with a given mass
type Sphere struct {
	Mass
}

// density returns the density of the sphere's material
func (s Sphere) density() float64 {
	return s.Density
}

// volume calculates the volume of the sphere given its radius
func (s Sphere) volume(radius float64) float64 {
	return (4.0 / 3.0) * math.Pi * math.Pow(radius, 3)
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////
// Cube represents a cube with a given mass
type Cube struct {
	Mass
}

// density returns the density of the cube's material
func (c Cube) density() float64 {
	return c.Density
}

// volume calculates the volume of the cube given its side length
func (c Cube) volume(side float64) float64 {
	return math.Pow(side, 3)
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Handler(massVolume MassVolume) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dimension, err := strconv.ParseFloat(r.URL.Query().Get("dimension"), 64); err == nil {
			weight := massVolume.density() * massVolume.volume(dimension)
			w.Write([]byte(fmt.Sprintf("Density : %f", massVolume.density())))
			w.Write([]byte(fmt.Sprintf("\nVolume : %f", massVolume.volume(dimension))))
			w.Write([]byte(fmt.Sprintf("\nTotal mass (Density*Volume) : "+"%.2f", math.Round(weight*100)/100)))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}
}


func WelcomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {  
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")  
		if _, err := w.Write([]byte("Hello TeamViewer")); err != nil {  
			http.Error(w, "Unable to write response", http.StatusInternalServerError)  
		}  
	}  
}


func main() {
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	aluminiumSphere := Sphere{Mass{Density: 2.710}}
	ironCube := Cube{Mass{Density: 7.874}}

	http.HandleFunc("/", WelcomeHandler())
	http.HandleFunc("/aluminium/sphere", Handler(aluminiumSphere))
	http.HandleFunc("/iron/cube", Handler(ironCube))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
