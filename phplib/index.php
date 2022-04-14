<?php

$localURL = "http://127.0.0.1:8000";
echo "<h2>test jwt auth</h2>";

echo "<hr><div><a
href='http://localhost/login?".
"&redirect_to=".$localURL."'>Auth</a></div>"; // auth link to accounts.bwretail.com

if ($_GET) {
    $jwtToken = $_GET['access_token'];
    $jwtArr = array_combine(['header', 'payload', 'signature'], explode('.', $jwtToken));
    echo "<pre>";
    echo "<pre>";
    $payload = base64_decode($jwtArr['payload']);
    $user = json_decode($payload);
    print_r($user); // token payload
    echo "</pre>";
    echo "<hr>";
    echo "<pre>";
    print_r($user);
    echo "<hr>";
    echo "<p>User: ". $user->eml . "</p>";
    echo "<p>Expiration: ". gmdate(DateTime::ISO8601, $user->exp) . "</p>";
}
