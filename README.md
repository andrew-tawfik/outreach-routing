# Outreach Routing Planner

## Project Description
This project is a Go-based Go-based desktop application designed to support Anba Abraam Service in planning efficient transportation routes for delivering groceries and transporting guests from dinner events to their homes. Built with Fyne UI framework, the application provides an intuitive drag-and-drop interface for coordinators to manage route dispatching with real-time visualization and interactive guest management.

The application integrates with Google Sheets for guest and event data, leverages external services like Google Maps API for geocoding and OSRM for calculating travel distances, and provides interactive map visualization for route planning.


At its core, the tool solves the Vehicle Routing Problem (VRP) using different optimization algorithms based on event type: Clarke-Wright Savings Algorithm for dinner events focused on optimal distance-based routing, and K-means++ clustering for grocery runs optimized for efficient geographic distribution. The GUI allows coordinators to review and manually adjust the automatically generated routes through an intuitive drag-and-drop interface.

---

## Key Features

- Interactive GUI Interface: Fyne-based application with drag-and-drop guest management between vehicles
- Real-time Map Visualization: Interactive Google Maps display showing vehicle routes and destinations
- State Management: Save, reset, and submit route changes with full undo capability
- Routes are displayed as vehicles and their assigned Guests
- Multi-tab Workflow: Organized interface with Home, Route Planning, and Map visualization tabs
- Reads and filters structured guest data from a Google Sheet
---

## Project Structure
```
/cmd/             → GUI application entry point using Fyne framework
/internal/
  ├── app/        → Core domain logic: routing algorithms, optimization, data models
  ├── converter/  → Data transformation between application layers
  ├── coordinates/→ Geographic coordinate handling and utilities
  ├── database/   → Google Sheets integration and data parsing
  ├── geoapi/     → Geocoding (Google Maps) and routing (OSRM) API integration
  └── ui/         → Fyne-based GUI components and user interaction 
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
