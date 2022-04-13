<?php

$localURL = "http://127.0.0.1:8000";
echo "<h2>test jwt auth</h2>";

echo "<hr><div><a
href='http://localhost/login?".
"&redirect_to=".$localURL."'>Auth</a></div>"; // auth link to accounts.bwretail.com

if (isset($_GET["code"])) {
    echo "Code: " . $_GET["code"]; //result of auth - access code
} else {
    echo "Code of access does not set";
}

echo "</pre>";
echo "<hr>";

if (isset($_GET["code"])) {
    $url =  "https://localhost/origin?"."code=".
    $_GET["code"]
    ."&application_id=a379ed35-a8e0-48c1-bfce-dc5eed92239c"; //get tokens by code
    $token_replay = file_get_contents($url);
    $repl = json_decode($token_replay);
    $tokens = $repl->Response;
}
if (isset($_GET["code"])) {
    echo "<pre>";
    var_dump($tokens->access_token);
    echo "</pre>";
    echo "<hr>";

    $jwtToken = $tokens->access_token;

    $jwtArr = array_combine(['header', 'payload', 'signature'], explode('.', $jwtToken));

    echo "<pre>";
    var_dump($jwtArr);
    echo "</pre>";
    echo "<hr>";
    echo "<pre>";
    print_r(base64_decode($jwtArr['header']));
    echo "</pre>";
    echo "<hr>";
    echo "<pre>";
    $payload = base64_decode($jwtArr['payload']);
    $user = json_decode($payload);
    print_r($user); // token payload
    echo "</pre>";
    echo "<hr>";
    echo "<p>User: ". $user->eml . "</p>";
    echo "<p>Expiration: ". gmdate(DateTime::ISO8601, $user->exp) . "</p>";
}
