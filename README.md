Ouvrir deux terminaux

## Dans /racine :
(Commandes de base dispo pour api et mongo)
Build dev (ou prod) + mongo	  ->  make build
Run les deux                  ->  make up
Stop les deux                 ->  make stop
Stop et supp les deux         ->  make down
Build & run les deux	        ->  make re
Montre les logs               ->  make logs
Terminal du docker	          ->  make bash-api
Terminal de mongo	            ->  make mongo-shell

## Dans /frontend :
python3 server.py

## Test Curl :
curl -X POST http://localhost:8080/upload -F "userName=Test" -F "codeFile=@test.zip"

curl -X POST http://localhost:8080/analyse -F "userName=Test"