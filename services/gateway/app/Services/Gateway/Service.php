<?php

declare(strict_types=1);

namespace App\Services\Gateway;

use App\Services\Gateway\Contracts\GatewayRepositoryContract;
use App\Services\Gateway\Events\RecordCreated;
use App\Services\Gateway\Contracts\ServiceContract;
use App\Services\Gateway\Models\Obit;
use Exception;

class Service implements ServiceContract {

    protected GatewayRepositoryContract $repository;

    /**
     * Service constructor.
     * @param GatewayRepositoryContract $repository
     */
    public function __construct(GatewayRepositoryContract $repository) {
        $this->repository = $repository;
    }

    /**
     * @param array $args
     * @return \Illuminate\Pagination\LengthAwarePaginator
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
     * @param string $obitDID
     * @return Obit
     */
    public function show(string $obitDID): ?Obit
    {
        return $this->repository->find($obitDID);
    }

    public function delete(string $obitId)
    {
        $obit = $this->update($obitId);
    }

    /**
     * @param string $obitDID
     * @return \Illuminate\Database\Eloquent\Relations\MorphMany
     * @throws Exception
     */
    public function history(string $obitDID)
    {
        $obit = $this->repository->find($obitDID);

        if (! $obit) {
            throw new Exception("Can't fetch the history because obit \"{$obitDID}\" not exists.");
        }

        return $obit->audits();
    }
}