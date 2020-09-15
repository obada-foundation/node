<?php

declare(strict_types=1);

namespace App\Services\Gateway;

use App\Services\Gateway\Contracts\GatewayRepositoryContract;
use App\Services\Gateway\Events\RecordCreated;
use App\Services\Gateway\Contracts\ServiceContract;
use App\Services\Gateway\Models\Obit;

class Service implements ServiceContract {

    /**
     * @var GatewayRepositoryContract
     */
    protected $repository;

    /**
     * Service constructor.
     * @param GatewayRepositoryContract $repository
     */
    public function __construct(GatewayRepositoryContract $repository) {
        $this->repository = $repository;
    }

    /**
     * @param array $args
     * @return \Illuminate\Support\Collection
     */
    public function search(array $args = []) {
        return $this->repository->findBy($args);
    }

    /**
     * @param ObitDto $dto
     * @return Obit
     */
    public function create(ObitDto $dto): Obit {
        $obit = new Obit;
        $obit->obit_did           = $dto->obitDID;
        $obit->usn                = $dto->usn;
        $obit->owner_did          = '';
        $obit->obd_did            = '';
        $obit->manufacturer       = '';
        $obit->part_number        = '';
        $obit->serial_number_hash = '';
        $obit->modified_at        = $dto->modifiedAt;
        $obit->root_hash          = $dto->rootHash();
        $obit->save();

        event(new RecordCreated($obit));

        return $obit;
    }

    public function update(string $obitId, ObitDto $dto)
    {
        // TODO: Implement update() method.
    }

    /**
     * @param string $obitId
     * @return Obit
     */
    public function show(string $obitId): ?Obit
    {
        return $this->repository->find($obitId);
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
