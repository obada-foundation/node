<?php

declare(strict_types=1);

namespace App\Services\Blockchain;

use App\Services\Blockchain\QLDB\Driver;
use App\Services\Blockchain\Events\RecordCreated;
use Throwable;
use Illuminate\Support\Facades\Log;

class Service implements ServiceContract {

    protected Driver $driver;

    /**
     * Service constructor.
     * @param Driver $driver
     */
    public function __construct(Driver $driver) {
        $this->driver = $driver;
    }

    /**
     * @param array $obit
     */
    public function create(array $obit) {
        try {
            $qldbObit = $this->driver->create($obit);
        } catch (Throwable $t) {
            Log::error(
                "Cannot submit obit: {$obit['obit_did']} to QLDB",
                [
                    'obit'      => $obit,
                    'exception' => $t->getTraceAsString()
                ]
            );

            throw $t;
        }

        event(new RecordCreated($qldbObit));
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
    public function history(string $obitId) {
        // TODO: Implement history() method.
    }

    /**
     * @param string $obitId
     * @param array $obit
     */
    public function update(string $obitId, array $obit) {
        try {
            $qldbObit = $this->driver->update($obitId, $obit);
        } catch (Throwable $t) {
            Log::error(
                "Cannot update obit: {$obitId} to QLDB",
                [
                    'obit'      => $obit,
                    'exception' => $t->getTraceAsString()
                ]
            );

            throw $t;
        }
    }

    /**
     * @param string $obitId
     * @return mixed
     */
    public function delete(string $obitId) {
        // TODO: Implement delete() method.
    }
}
