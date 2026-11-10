# API Telemetria Robo

API para o recebimento de dados dos sensores do robô e envio ao front-end para analise posterior.

## Como executar

No diretório do código, executar o comando abaixo para iniciar a API.

```
go run main.go
```

## Endpoints

- /post-match

Métodos aceitos: POST

Exemplo corpo da requisição:
```
{
    "title": "MatchTitle",
    "date": "1970-01-01",
    "opponentName": "AdversaryName"
}
```

Exemplo corpo de resposta:
```
{
    "status": "success"
}
```

- /get-closed-match-id/{id:int}

Métodos aceitos: GET

Exemplo corpo da requisição:
```
(indiferente)
```

Exemplo corpo de resposta:
```
{
    "id": 1,
    "title": "MatchName",
    "date": "1970-01-01",
    "opponentName": "AdversaryName",
    "closed": true,
    "roundCount": 3
}
```

- /get-closed-matches

Métodos aceitos: GET

Exemplo corpo da requisição:
```
(indiferente)
```

Exemplo corpo de resposta:
```
[
    {
        "id": 1,
        "title": "MatchName",
        "date": "1970-01-01",
        "opponentName": "AdversaryName",
        "closed": true,
        "roundCount": 3
    },
    {
        "id": 2,
        "title": "MatchName2",
        "date": "1970-01-02",
        "opponentName": "Adversary2Name",
        "closed": true,
        "roundCount": 2
    }
]
```

- /get-current-match

Métodos aceitos: GET

Exemplo corpo da requisição:
```
(indiferente)
```

Exemplo corpo de resposta:
```
{
    "id": 1,
    "title": "MatchName",
    "date": "1970-01-01",
    "opponentName": "AdversaryName",
    "closed": false
}
```

- /close-current-match

Métodos aceitos: POST

Exemplo corpo da requisição:
```
(indiferente)
```

Exemplo corpo de resposta:
```
{
    "status": "success"
}
```

- /post-round

Métodos aceitos: POST

Exemplo corpo da requisição:
```
{
    "sensors": "(String Base64 de arquivo csv)"
}
```

Exemplo corpo de resposta:
```
{
    "status": "success"
}
```

- /get-match-round/{id:int}/{roundNumber:int}

Métodos aceitos: POST

Exemplo corpo da requisição:
```
(indiferente)
```

Exemplo corpo de resposta:
```
{
    "roundNumeber": 1,
    "sensors": [
        {
            "name": "RPM",
            "readings": [
                {
                    "timestamp": "00:000050",
                    "value": 1231
                },
                {
                    "timestamp": "00:000100",
                    "value": 1231
                }
            ]
        },
        {
            "name": "S1",
            "readings": [
                {
                    "timestamp": "00:000050",
                    "value": 0
                },
                {
                    "timestamp": "00:000100",
                    "value": 1
                }
            ]
        },
        {
            "name": "S2",
            "readings": [
                {
                    "timestamp": "00:000050",
                    "value": 1
                },
                {
                    "timestamp": "00:000100",
                    "value": 1
                }
            ]
        },
        {
            "name": "QTR1",
            "readings": [
                {
                    "timestamp": "00:000050",
                    "value": 0
                },
                {
                    "timestamp": "00:000100",
                    "value": 0
                }
            ]
        },
        {
            "name": "QTR2",
            "readings": [
                {
                    "timestamp": "00:000050",
                    "value": 0
                },
                {
                    "timestamp": "00:000100",
                    "value": 1
                }
            ]
        }
    ]
}
```