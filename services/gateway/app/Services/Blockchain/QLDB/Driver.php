<?php

declare(strict_types=1);

namespace App\Services\Blockchain\QLDB;

use App\Services\Blockchain\Ion;
use Aws\QLDBSession\QLDBSessionClient;

class Driver  {

    protected QLDBSessionClient $qldb;

    protected string $sessionToken;

    /**
     * Driver constructor.
     * @param QLDBSessionClient $qldb
     * @param Ion $ion
     */
    public function __construct(QLDBSessionClient $qldb) {
        $this->qldb = $qldb;

        $session = $this->qldb->sendCommand([
            'StartSession' => [
                'LedgerName' => config('qldb.ledger_name')
            ]
        ])->toArray();

        $this->sessionToken = $session['StartSession']['SessionToken'];
    }

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
