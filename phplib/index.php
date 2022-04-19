<?php

function base64url_encode($str)
{
    return rtrim(strtr(base64_encode($str), '+/', '-_'), '=');
}

function is_jwt_valid($jwt, $secret = 'secret')
{
    // split the jwt
    $tokenParts = explode('.', $jwt);
    $header = base64_decode($tokenParts[0]);
    $payload = base64_decode($tokenParts[1]);
    $signature_provided = $tokenParts[2];

    // check the expiration time - note this will cause an error if there is no 'exp' claim in the jwt
    $expiration = json_decode($payload)->exp;
    $is_token_expired = ($expiration - time()) < 0;

    // build a signature based on the header and payload using the secret
    $base64_url_header = base64url_encode($header);
    $base64_url_payload = base64url_encode($payload);
    $signature = hash_hmac('SHA256', $base64_url_header . "." . $base64_url_payload, $secret, true);
    $base64_url_signature = base64url_encode($signature);

    // verify it matches the signature provided in the jwt
    $is_signature_valid = ($base64_url_signature === $signature_provided);
    if ($is_token_expired || !$is_signature_valid) {
        return false;
    } else {
        return true;
    }
}

$localURL = "https://bwretail.dev/jwt/";
$app_id = "a379ed35-a8e0-48c1-bfce-dc5eed92239c";
$app_secret = "3dp9gudw0l19yr9ois8iu9b3220qemn8";
echo "<h2>test jwt auth</h2>";

echo "<hr><div><a
href='http://accounts.bwretail.com/login?&redirect_to={$localURL}&app_id={$app_id}'>Auth</a></div>"; // auth link to accounts.bwretail.com

if ($_GET) {
    $jwtToken = $_GET['access_token'];

    $res = is_jwt_valid($jwtToken, $app_secret);
    echo "<hr><div><p>";
    if ($res) {
        echo "<span style='color:blue'>JWT token is valid</span>" . PHP_EOL;
    } else {
        echo "<span style='color:red'>JWT token is not valid</span>" . PHP_EOL;
    }
    echo "</p></div><hr><hr>";

    $jwtArr = array_combine(['header', 'payload', 'signature'], explode('.', $jwtToken));
    echo "<pre>";
    $payload = base64_decode($jwtArr['payload']);
    $user = json_decode($payload);
    print_r($user); // token payload

    echo "</pre>";
    echo "<hr>";
    echo "<hr>";
    echo "<p>User: ". $user->eml . "</p>";
    echo "<p>Expiration: ". gmdate(DateTime::ISO8601, $user->exp) . "</p>";
}
