Meetup RSVP exporter
====================

``meetup-rsvpexport`` is a command-line tool to export RSVPs for a given
Meetup event to CSV format. It exports:

* Attendee name
* Attendee profile bio
* RSVP status (yes/no)
* Number of guests


Why this tool?
--------------

I needed a tool to export RSVPs for my events to print name tags on A4 paper.
Unfortunately Meetup only offers a way to print name tags on Letter sized
format.


How to install
--------------

Make sure you have Go installed and ``$GOPATH`` is set. Then you should
be able to install this tool with::

    $ go get github.com/brocaar/meetup-rsvpexport


Usage
-----

Before you are able to use this tool, you need to request API access.
See: http://www.meetup.com/meetup_api/

::

    Usage of ./meetup-rsvpexport:
      -apikey="": your meetup API key
      -eventid="": the event-id of the meetup

Example::

    $ meetup-rsvpexport -apikey yourapikey -eventid 12345 > attendees.csv
