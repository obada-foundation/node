<?php

declare(strict_types=1);

namespace App\Services\Gateway;

use App\Services\Gateway\Contracts\ServiceContract;
use App\Services\Gateway\Models\Obit;

class Service implements ServiceContract {

    public function __construct() {

    }

    public function search() {
        // TODO: Implement search() method.
    }

    /**
     * @param ObitDto $dto
     * @return mixed|void
     */
    public function create(ObitDto $dto) {
        $obit = new Obit;
        $obit->obit_did           = $dto->obitDID;
        $obit->usn                = $dto->usn;
        $obit->owner_did          = '';
        $obit->obd_did            = '';
        $obit->manufacturer       = '';
        $obit->part_number        = '';
        $obit->serial_number_hash = '';
        $obit->modified_at        = $dto->modifiedAt;
        $obit->save();

        dd($obit);
    }

    public function update(string $obitId)
    {
        // TODO: Implement update() method.
    }

    public function show(string $obitId)
    {
        // TODO: Implement show() method.
    }

    public function delete(string $obitId)
    {
        // TODO: Implement delete() method.
    }

    public function history(string $obitId)
    {
        // TODO: Implement history() method.
    }
}
