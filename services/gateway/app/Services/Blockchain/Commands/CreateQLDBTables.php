<?php

declare(strict_types=1);

namespace App\Services\Blockchain\Commands;

use Aws\QLDB\QLDBClient;
use Aws\QLDBSession\QLDBSessionClient;
use Illuminate\Console\Command;

class CreateQLDBTables extends Command {

    protected $signature = 'qldb:create-tables';

    protected QLDBClient $qldb;

    protected QLDBSessionClient $sessionClient;

    public function __construct(QLDBClient $qldb, QLDBSessionClient $sessionClient)
    {
        parent::__construct();

        $this->qldb = $qldb;
        $this->sessionClient = $sessionClient;
    }

    public function handle() {
        $this->qldb->createLedger([
            'DeletionProtection' => false,
            'Name'               => config('qldb.ledger_name'),
            'PermissionsMode'    => 'ALLOW_ALL'
        ]);

        sleep(40);

        $session = $this->sessionClient->sendCommand([
            'StartSession' => [
                'LedgerName' => config('qldb.ledger_name')
            ]
        ])->toArray();

        $sessionToken = $session['StartSession']['SessionToken'];

        $tx = $this->sessionClient->sendCommand([
            'SessionToken' => $sessionToken,
            'StartTransaction' => [

            ]
        ])->toArray();

        $transactionId = $tx['StartTransaction']['TransactionId'];

        $this->sessionClient->sendCommand([
            'SessionToken' => $sessionToken,
            'ExecuteStatement' => [
                'TransactionId' => $transactionId,
                'Statement' => 'CREATE TABLE Obits',
            ]
        ])->toArray();

        $r = $this->sessionClient->sendCommand([
            'SessionToken' => $sessionToken,
            'CommitTransaction' => [
                'TransactionId' => $transactionId
            ]
        ])->toArray();

        dd($r);

        foreach ($r['ExecuteStatement']['FirstPage']['Values'] as $x) {
            $arr = implode(" ", unpack("C*", $x['IonBinary']));
        }

        dd($arr);
    }
}
