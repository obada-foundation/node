<?php

declare(strict_types=1);

namespace App\Services\Blockchain;

use App\Services\Blockchain\Contracts\ServiceContract;
use Aws\QLDBSession\QLDBSessionClient;

class Service implements ServiceContract {

    /**
     * @var QLDBSessionClient
     */
    protected $qldb;

    /**
     * Service constructor.
     * @param QLDBSessionClient $qldb
     */
    public function __construct(QLDBSessionClient $qldb) {
        $this->qldb = $qldb;
    }

    /**
     * @return mixed
     */
    public function create() {
        $session = $this->qldb->sendCommand([
            'StartSession' => [
                'LedgerName' => config('qldb.ledger_name')
            ]
        ])->toArray();

        $sessionToken = $session['StartSession']['SessionToken'];

        $tx = $this->qldb->sendCommand([
            'SessionToken' => $sessionToken,
            'StartTransaction' => [

            ]
        ])->toArray();

        $transactionId = $tx['StartTransaction']['TransactionId'];

        $r = $this->qldb->sendCommand([
            'SessionToken' => $sessionToken,
            'ExecuteStatement' => [
                'TransactionId' => $transactionId,
                'Statement' => 'SELECT * FROM DriversLicense'
            ]
        ])->toArray();


        foreach ($r['ExecuteStatement']['FirstPage']['Values'] as $x) {
            echo "\n\n";
            $arr = unpack("C*", $x['IonBinary']);
            echo implode(" ", $arr);
        }
    }

    /**
     * @param string $obitId
     * @return mixed
     */
    public function show(string $obitId)
    {
        // TODO: Implement show() method.
    }

    /**
     * @param string $obitId
     * @return mixed
     */
    public function history(string $obitId)
    {
        // TODO: Implement history() method.
    }

    /**
     * @param string $obitId
     * @return mixed
     */
    public function update(string $obitId)
    {
        // TODO: Implement update() method.
    }

    /**
     * @param string $obitId
     * @return mixed
     */
    public function delete(string $obitId)
    {
        // TODO: Implement delete() method.
    }
}
