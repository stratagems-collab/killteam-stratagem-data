# killteam-stratagem-data

Updates from KT Annual '19 included.
Please contribute by *reporting issues* using github issues, or *creating pull-requests*.

Format documentation
```
{
	"packageName": "KillTeam Stratagems",
	"versionName": "0.1-alpha",
	"versionCode": 1,
	"factions": [
		{
			"faction": "Name of faction, used for matching",
			"tactics": [
				{
					//required
					"title": "Title",
					"sub": "Faction Name Tactic",
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
						"Keyword",
						"optional"
					],
					"equipment": [
						"weapons",
						"equipment",
						"optional"
					],
					"specialist": "Sniper, or other specialism this is specific for",
					"level": 1,
					"source": {
						"id": "id_of_source",
						"page": "page",
						"data": "data for url"
					}
				}
			]
		}
	],
	"sources": [
		{
			"id": "id of source, for referencing from tactics",
			"title": "Title",
			"url": "http://example.com?ref=%s"
		}
	]
}
```
