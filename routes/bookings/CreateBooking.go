package bookings

import (
	"carpool-backend/database"
	"carpool-backend/helpers"
	"carpool-backend/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

func CreateBooking(c *fiber.Ctx) error {
	// Define a variable to hold the Booking struct
	var Booking models.Booking

	// Parse the request body into the Booking struct
	err := c.BodyParser(&Booking)
	if err != nil {
		log.Printf("Error parsing JSON: %v\n", err)
		return c.Status(400).SendString("Error parsing JSON")
	}

	// Validate the Booking struct
	if err := helpers.ValidateBooking(Booking); err != nil {
		log.Printf("Error validating Booking: %v\n", err)
		return c.Status(400).SendString("Error validating Booking")
	}

	// Check if the user has already booked the ride
	existingBooking := models.Booking{}
	if database.Database.Db.Where("user_id = ? AND ride_id = ?", Booking.PassengerID, Booking.RideID).First(&existingBooking).Error == nil {
		log.Println("Similar booking already exists")
		return c.Status(400).SendString("Similar booking already exists")
	}
	ride := models.Ride{}
	if err := database.Database.Db.Where("id = ?", Booking.RideID).First(&ride).Error; err != nil {
		log.Printf("Error finding ride: %v\n", err)
		return c.Status(404).SendString("Ride not found")
	}

	// Check if there are available seats
	if ride.BookedSeats >= ride.TotalSeats {
		log.Println("No available seats for this ride")
		return c.Status(400).SendString("No available seats for this ride")
	}

	// Update the booked seats count
	ride.BookedSeats++
	if err := database.Database.Db.Save(&ride).Error; err != nil {
		log.Printf("Error updating booked seats for ride: %v\n", err)
		return c.Status(500).SendString("Error updating booked seats for ride")
	}

	// Log the updated booked seats for the ride
	log.Printf("Booked seats for ride with id %v incremented to %v\n", ride.ID, ride.BookedSeats)
	// Insert the booking into the database
	result := database.Database.Db.Create(&Booking)
	if result.Error != nil {
		log.Printf("Error creating booking: %v\n", result.Error)
		return c.Status(500).SendString("Error creating booking")
	}

	// Log the booking id
	log.Printf("Booking with id %v created\n", Booking.ID)
	return c.Status(200).JSON(Booking)
}
