# EEBUS

This is an open source implementation of parts of the [EEBUS protocol][1] Version 1.0.1 specification in [Go][2]

**WARNING** The code in this repository is a mess and we know it!

## Current State

- This should be considered as a proof of concept
- The main goal (for now) was to get a working implementation for the EV use cases
- The implementation does not yet have a clean and easy to understand architecture
- The available documentation and missing reference implementation lead to the current state and will evolve from there
- The main purpose of this implementation is being able to access EEBUS compatible chargers in [EVCC][3]
- Only EV specific use cases are being worked on right now

## Missing

- Adopt a proper code architecture
- Proper error handling
- Code cleanups
- Increasing test coverage
- Documentation
- Proper APIs for public use

## Features

- Partly implemented EEBUS Use Cases:
  - EVSE Commissioning and Configuration
  - EV Commissioning and Configuration
  - EV Charging Electricity Measurement
  - EV State of Charge
  - Optimization of Self-Consumption During EV Charging
  - Overload Protection by EV Charging Current Curtailment

## Background

Having an EV with a bundled charger that basically only provides an EEBUS interface, the wish was to be able to control charging of the EV via PV. [EVCC][3] provided already a great foundation, API for different chargers, PV charging support, but had no way to connect with chargers via the EEBUS protocol. Even though the specification of EEBUS are open source, there seems to be no non commercial implementation existing. So [andig](https://github.com/andig/) and [Andreas Linde](https://github.com/DerAndereAndi) took up the challenge and see where we can get.

[1]: https://www.eebus.org
[2]: https://golang.org
[3]: https://evcc.io
