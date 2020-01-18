# killteam-stratagem-data

Format documentation
```
{
	"packageName": "KillTeam Stratagems",
	"versionName": "0.1-alpha",
	"versionCode": 1,
	"factions": [
		{
			"faction": "Name of faction, used for matching",
			"Stratagems": [
				{
					"title": "Title of Tactic, required",
					"sub": "Faction Name Tactic, required",
					"desc": "The description on how this tactic works. required",
					"cp": 2,
					"phases": {
							"move": true,
							"psychic": false,
							"shoot": false,
							"fight": false,
							"round": false,
							"round_one": false,
							"morale": false,
							"event": false
					},
					"keywords": [
						"array of Keywords, providing keywords is optional"
					],
					"equipment": [
						"array of weapons or equipment, optional"
					],
					"specialist": "Sniper, or other specialism this is specific for",
					"level": 1
				}
			]
		}
	]
}
```
