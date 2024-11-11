# Linkanything

Provides news streams about anything.

## Endpoints

### Create Link
```
POST https://{{HOST}}/links HTTP/1.1
```

```json
{
    "title": "Starboard Wingboard Foil 2024",
    "uri": "https://www.wingdaily.de/news-starboard-wingboard-foil-2024-20230928.htm",
    "tags": ["Test"]
}
```

### Get Links
```
GET https://{{HOST}}/links[?page=1&limit=20]
```
```json
{
  "links": [
    {
      "id": 24,
      "title": "Guter Erfahrungsbericht / Thread zur Axis Spitfire Range",
      "uri": "https://oaseforum.de/showthread.php?t=184756",
      "tags": [
        "Report"
      ],
      "created": "2023-10-01T13:57:47.210661Z"
    },
    {
      "id": 22,
      "title": "Starboard Wingboard Foil 2024",
      "uri": "https://www.wingdaily.de/news-starboard-wingboard-foil-2024-20230928.htm",
      "tags": [
        "Test"
      ],
      "created": "2023-09-29T08:10:38.337925Z"
    }
  ],
  "pagination": {
    "total_records": 2,
    "current_page": 1,
    "total_pages": 10,
    "next_page": 2,
    "prev_page": null
  }
}
```

## Development

Install sqlc

```
brew install sqlc
```