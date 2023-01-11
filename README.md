# UKSFCG
UK Sector File Changelog Generator is a small CLI tool designed to make releasing the UK Sector File easier by separating the kinds of changes and contributors, to be easily copied and pasted.

The output is in this format:
```txt
--- AIRACs: ---
2207:
Updated Dundee (EGPN) hold D on SMR
Updated Cranfield (EGTC) SMR
Updated Inverness (EGPE) RWYs 05/23 and 11/29 coords
Removed Belfast Aldergrove (OY) NDB

--- Other: ---
Bug - Fixed D087E Position
Enhancement - Updated Belfast Aldergrove (EGAA) SMR to better differentiate Tug and Hold Points
Bug - Corrected Alderney (EGJA) runway coords
Enhancement - Added missing heli points and holds to Gloucestershire (EGBJ) SMR
Enhancement - Refined Fairford (EGVA) SMR
Enhancement - Added Luton (EGGW) Tug Release Points
Enhancement - Added Stansted (EGSS) Ground Network
Bug - Corrected Guernsey (EGJB) runway threshold coords
Enhancement - Added Gatwick (EGKK) Ground Network
Bug - Corrected Prestwick (EGPK) RWY 30 threshold coords
Enhancement - Added Luton (EGGW) Ground Network
Enhancement - Added Glasgow (EGPF) Ground Network
Enhancement - Added Liverpool (EGGP) Ground Network

--- Contributors: ---
John Doe
Smith John
Larry Benwater
```

## Usage
The easiest way to use this is copy and paste the desired changelog list from [https://raw.githubusercontent.com/VATSIM-UK/UK-Sector-File/main/.github/CHANGELOG.md](https://raw.githubusercontent.com/VATSIM-UK/UK-Sector-File/main/.github/CHANGELOG.md) into a `changelog.md` file locally, and just run with `./UKSFCG`

Flags:
- `--in` - sets the input file, optional, default `changelog.md`
- `--out` - sets the output file, optional, default `output.txt`

## Limitations
- If a contributor has been added in the format `- thanks to @John (John Smith) and Smith John`, Smith John would not be on the contributors list
