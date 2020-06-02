<?php

namespace GuzzleClient;

require __DIR__.'/vendor/autoload.php';

use GuzzleHttp\Client;

$client = new Client([
    "headers" => [
        'Accept' => 'application/json',
        'Content-Type' => 'application/json'
    ]
]);

echo "Bank accounts : \n";
$response = $client->get("https://lpcloud-project-1.ew.r.appspot.com/");

echo $response->getBody()->getContents();

echo "\n Adding a bank account \n";

$bankAccount = [
    "montant" => 250000,
    "nom" => "Test",
    "prenom" => "Test",
    "risk" => "low"
];

$jsonEncodedBA = json_encode($bankAccount);
$responsePost = $client->put(
    "https://lpcloud-project-1.ew.r.appspot.com/",
    ["body" => $jsonEncodedBA]
);

echo "\n New list :\n";

$response = $client->get("https://lpcloud-project-1.ew.r.appspot.com/");

echo $response->getBody()->getContents();


echo "\n Sending a loan of 8000€ for the test client\n";

$response = $client->get("https://lpro-cloud-loan-approval-1.herokuapp.com/?nom=Test&montant=8000");

echo $response->getBody()->getContents();

echo "\n listing Accounts: \n";

$response = $client->get("https://lpcloud-project-1.ew.r.appspot.com/");

echo $response->getBody()->getContents();

echo "\n listing Approvals : \n";

$response = $client->get("https://lp-cloud-app-manager-1.ew.r.appspot.com/");

echo $response->getBody()->getContents();

echo "\n Sending a loan of 12000€ for the test client \n";

$response = $client->get("https://lpro-cloud-loan-approval-1.herokuapp.com/?nom=Test&montant=12000");

echo $response->getBody()->getContents();

echo "\n Listing Approvals : \n";

$response = $client->get("https://lp-cloud-app-manager-1.ew.r.appspot.com/");

echo $response->getBody()->getContents();