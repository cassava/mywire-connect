mywire-connect
==============

Die Firma [mywire Datentechnik GmbH](my-wire.de) stellt unter anderem
Internetzugang für Studenten in bestimmte Studentenwohnheime zur Verfügung.
Nervig ist, dass man sich über die mywire Webseite sich anmelden muss, bevor
man Internetzugang hat.

Das Programm **mywire-connect** meldet einem automatisch an, wenn eine
mywire Anmeldung erforderlich ist. Dabei ist es egal, an welchem (von
mywire vermutlich geliefertem) Modem man den Computer angeschlossen hat.

Voraussetzungen
---------------
Die Umgebungsvariablen `MYWIRE_USER` und `MYWIRE_PASS` müssen gesetzt sein,
sonst kann mywire-connect eine Anmeldung nicht ausführen.

Installation
------------
mywire-connect ist in die Programmiersprache [Go](golang.org) geschrieben.
Kompilation erzeugt ein statisch gelinktes Programm, was man überall hintun
kann und dort ausführen.

Sie können binäre Dateien für die aktuelle Version 0.9 hier
[runterladen](https://github.com/cassava/mywire-connect/releases/tag/v0.9).


Falls Sie selber kompilieren wollen, müssen Sie Go zuerst installieren.
Danach können sie mit dem Befehl `go get github.com/cassava/mywire-connect`
das Programm runterladen sowie kompilieren. Es müsste eine ausführbare Datei
dann vorhanden sein.
