<?php

declare(strict_types=1);

namespace App\Services\Gateway;

use App\Services\Gateway\Repositories\GatewayRepositoryContract;
use App\Services\Gateway\Events\RecordCreated;
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
     * @param ObitDto $dto
     * @return string
     */
    public function buildRootHash(ObitDto $dto): string {
        $lastObit = Obit::orderBy('id', 'DESC')->first();

        $parentRootHash = $lastObit ? $lastObit->root_hash : null;

        $hashedData = hash('sha256', sprintf(
            '%s%s%s%s%s%s',
            $dto->obitDID,
            $dto->usn,
            $dto->manufacturer,
            $dto->partNumber,
            $dto->serialNumberHash,
            $dto->modifiedAt
        ));

        return hash('sha256', $parentRootHash . $hashedData);
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
        $obit->parent_id          = optional(Obit::orderBy('id', 'DESC')->first())->id;
        $obit->obit_did           = $dto->obitDID;
        $obit->usn                = $dto->usn;
        $obit->owner_did          = '';
        $obit->obd_did            = '';
        $obit->manufacturer       = $dto->manufacturer;
        $obit->part_number        = $dto->partNumber;
        $obit->serial_number_hash = $dto->serialNumberHash;
        $obit->modified_at        = $dto->modifiedAt;
        $obit->root_hash          = $this->buildRootHash($dto);
        $obit->save();

        event(new RecordCreated($obit));

        return $obit;
    }

    public function update(string $obitId, ObitDto $dto) {
        // TODO: Implement update() method.
    }

    /**
     * @param string $obitDID
     * @return Obit
     */
    public function show(string $obitDID): ?Obit {
        return $this->repository->find($obitDID);
    }

    public function delete(string $obitId) {
        $obit = $this->update($obitId);
    }

    /**
     * @param string $obitDID
     * @return \Illuminate\Database\Eloquent\Relations\MorphMany
     * @throws Exception
     */
    public function history(string $obitDID) {
        $obit = $this->repository->find($obitDID);

        if (! $obit) {
            throw new Exception("Can't fetch the history because obit \"{$obitDID}\" not exists.");
        }

        return $obit->audits();
    }
}
