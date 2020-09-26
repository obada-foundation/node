<?php

declare(strict_types=1);

namespace App\Services\Blockchain\QLDB;

use App\Services\Blockchain\Ion;
use Aws\QLDBSession\QLDBSessionClient;

class Driver  {

    protected QLDBSessionClient $qldb;

    protected string $sessionToken;

    public function create(array $obit) {
        $tx = $this->qldb->sendCommand([
            'SessionToken' => $this->sessionToken,
            'StartTransaction' => [

            ]
        ])->toArray();

        unset($obit['id']);

        $transactionId = $tx['StartTransaction']['TransactionId'];

        $r = $this->qldb->sendCommand([
            'SessionToken' => $this->sessionToken,
            'ExecuteStatement' => [
                'TransactionId' => $transactionId,
                'Statement' => 'INSERT INTO Obits ?',
                'Parameters' => [
                    'IonText' => $obit
                ]
            ]
        ])->toArray();
    }
}
