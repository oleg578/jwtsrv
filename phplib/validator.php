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

$jwt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWwiOiJvbGVnLm5hZ29ybmlqQGdtYWlsLmNvbSIsImV4cCI6MTY1MDI5NjU3MCwianRpIjoiMzQzMjY2NjYtNjI2Mi00OTM5LWFkNjEtMzUzMjYzMmQzNDM5IiwibmljayI6Ik9sZWgiLCJyb2xlIjoiZ3Vlc3QiLCJ1aWQiOiI0MmZmYmI5OS1hNTJjLTQ5YmEtODhlNC00NzU1YjA4MWNhYTYifQ.tTl98gkhYYuyyOC44rQnLe36PFKp8RtGKehJfrw0dX8";

$secret = "3dp9gudw0l19yr9ois8iu9b3220qemn8";

$res = is_jwt_valid($jwt, $secret);
if ($res) {
    echo "Valid JWT" . PHP_EOL;
} else {
    echo "Invalid JWT" . PHP_EOL;
}
