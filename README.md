# UKSFCG
UK Sector File Changelog Generator is a small CLI tool designed to make releasing the UK Sector File easier by separating the kinds of changes and contributors, to be easily copied and pasted.

The output is in this format:

```txt
--- AIRACs: ---
2309:
Updated Kirknewton (EGKT) Runway Coordinate
Updated Dover (DVR) VOR/DME coordinates
Added EG D713 (Fast Jet Area South) & EG D901 (Fast Jet Area North)
Added Portland Heliport (EGDP) A/G Position and Updated Basic Data
Updated ADN, IOM, LND and TNT VOR Coordinates
Removed Bristol (EGGD) stands 31R & 33L
Added Waddignton (EGXW) RNP approach fixes
Updated Wattisham (EGUW) Tower frequency
Added Biggin Hill (EGKR) new taxiway
Added Fairoaks (EGTF) disused taxiway

2308:
Updated Belfast Aldergrove (EGAA) SMR

--- Other: ---
Procedure Change (2309):
Amended EuroCenter vACC EURM-W and EURW-N callsigns to EUC-MW and EUC-WN, respectively
8.33KHz Trial - Changed non-UK frequencies
8.33KHz Trial (ENR Phase 1) - Transitioned LAC West & Clacton frequencies
8.33 Trial (AD Phase 1) - Transition EGLL/PH/SS/GP Frequencies

Bug:
Removed Gatwick (EGKK) stand 145 L/R labels
Renamed Waddington (EGXW) holding points
Fixed Farnborough (EGLF) Inbound Agreements from the north via CPT.

Enhancement:
Updated Birmingham (EGBB) SMR style
Added Derby (EGBD) SMR
Removed Manchester (EGCC) disused stands
Improved East Midlands (EGNX) VFR Lane and SID Line Colours
Added Jersey Control radar regions & colours
Added Yeovilton (EGDY) stand numbers
Added Leeds East (EGCM) SMR
Updated Barton (EGCB) SMR
Added Haverfordwest (EGFE) SMR

--- Contributors: ---
John Doe
Smith John
Larry Benwater
```

## Usage
There are two ways to use this program:
- Locally - copy and paste the desired changelog list from [https://raw.githubusercontent.com/VATSIM-UK/UK-Sector-File/main/.github/CHANGELOG.md](https://raw.githubusercontent.com/VATSIM-UK/UK-Sector-File/main/.github/CHANGELOG.md) into a `changelog.md` file locally, and just run with `./UKSFCG`
- Online - run the program using `./UKSFCG --url`, which will read from the default URL specified below

Flags:
- `--in` - sets the input file, optional, default `changelog.md`
- `--out` - sets the output file, optional, default `output.txt`
- `--url` - sets the url, optional, if `--url` is present but empty then it defaults to `https://raw.githubusercontent.com/VATSIM-UK/UK-Sector-File/main/.github/CHANGELOG.md` 

The `--url` flag takes precedence over `--in`, so if both are specified, then it will read from online.

## Limitations
- If a contributor has been added in the format `- thanks to @John (John Smith) and Smith John`, Smith John would not be on the contributors list
