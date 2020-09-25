<?php

return [
    'connection' => [
        'credentials' => [
            'key'     => env('AWS_KEY'),
            'secret'  => env('AWS_SECRET'),
        ],
        'region'  => env('AWS_REGION', 'us-east-1'),
        'version' => 'latest'
    ],
    'ledger_name' => env('QLDB_LEDGER_NAME', 'obada')
];
