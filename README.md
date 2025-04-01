# Outreach Routing Planner

## Project Description
This project is a Go-based command-line tool built to support Anba Abraam Service in planning efficient transportation routes for delivering groceries and transporting guests to dinner events. Designed specifically for the Service's coordinators, it streamlines the entire route dispatching process.

The tool integrates with Google Sheets for guest and event data, and leverages external services like Nominatim for geocoding and OSRM for calculating travel distances. 

At its core, the tool solves the Vehicle Routing Problem (VRP). This is a classic optimization problem in Computer Science that involves determining the most efficient set of routes for a fleet of vehicles to serve a group of destinations. In this case, destinations are guest dropoff points, and the fleet consists of vehicles with limited seating capacity. The system uses the Clarke-Wright Savings Algorithm to group guests intelligently and minimize total travel distance while ensuring each vehicle's capacity is respected.

---

## Key Features

- Reads structured guest and event data from a Google Sheet
- Filters guests based on service eligibility
- Retreive exact coordinates using guest addresses using Nominatim API
- Builds a distance matrix of all coordinates involved using OSRM (Open Source Routing Machine)
- Determines optimised vehicle routes using Clarke-Wright Savings Algorithm
- Routes are displayed as vehicles and their assigned Guests
---

## Project Structure
```
/cmd/             → CLI entrypoint and orchestration
/internal/
  └── database/   → Fetches and parses Google Sheets data
  └── geoapi/     → Handles geocoding (Nominatim) and routing (OSRM) APIs
  └── app/        → Core domain logic: routing, distance handling, optimization
```


## Input Format

### Headers
Google Sheet must have the following column headers:
```
Status | Name | Group Size | Number | Address
```
### Address Formatting Guidelines
To ensure accurate geocoding and routing:

- All guest addresses must be valid Ottawa addresses.

- Only "unit" or "apt" are allowed as suffixes to the address. These are stripped automatically during processing.

  - Example of acceptable input: 96 George St apt 54

- Avoid extra text or notes in the address field, as this can lead to geocoding failures.
