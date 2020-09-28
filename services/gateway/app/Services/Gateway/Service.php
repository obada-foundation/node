<?php

declare(strict_types=1);

namespace App\Services\Gateway;

use App\Services\Gateway\Repositories\GatewayRepositoryContract;
use App\Services\Gateway\Events\RecordCreated;
use App\Services\Gateway\Models\Obit;
use Exception;
use Illuminate\Support\Collection;
use Illuminate\Support\Facades\Log;
use Throwable;

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
        $obit->obit_status        = $dto->obitStatus;
        $obit->owner_did          = (string) $dto->ownerDID;
        $obit->obd_did            = (string) $dto->obdDID;
        $obit->metadata           = $dto->metadata;
        $obit->doc_links          = $dto->docLinks;
        $obit->structured_data    = $dto->structuredData;
        $obit->manufacturer       = $dto->manufacturer;
        $obit->part_number        = $dto->partNumber;
        $obit->serial_number_hash = $dto->serialNumberHash;
        $obit->modified_at        = $dto->modifiedAt;
        $obit->root_hash          = $this->buildRootHash($dto);
        $obit->save();

        event(new RecordCreated($obit));

        return $obit;
    }

    /**
     * @param string $obitId
     * @param UpdateObitDto $dto
     * @return mixed|void
     */
    public function update(string $obitId, UpdateObitDto $dto) {
        $obit = $this->repository->find($obitId);

        $update = collect($dto->toArray())
            ->filter(fn ($v) => $v != null)
            ->flip()
            ->map(fn ($v) => strtolower(preg_replace(['/([a-z\d])([A-Z])/', '/([^_])([A-Z][a-z])/'], '$1_$2', $v)))
            ->flip()
            ->toArray();

        $obit->update($update);
    }

    /**
     * @param string $obitDID
     * @return Obit
     */
    public function show(string $obitDID): ?Obit {
        return $this->repository->find($obitDID);
    }

    public function delete(string $obitId) {
        $obit = $this->repository->find($obitId);
        $obit->obit_status = Obit::DISABLED_BY_OWNER_STATUS;
        $obit->save();
    }

    /**
     * @param string $obitDID
     * @return Collection
     * @throws Exception
     */
    public function history(string $obitDID): Collection {
        $obit = $this->repository->find($obitDID);

        if (! $obit) {
            throw new Exception("Can't fetch the history because obit \"{$obitDID}\" not exists.");
        }

        return $obit->audits;
    }

    /**
     * @param string $obitDID
     * @throws Throwable
     */
    public function commit(string $obitDID) {
        $obit = $this->repository->find($obitDID);

        if (! $obit) {
            throw new Exception("Cannot commit. Given obit: \"{$obitDID}\" not exists.");
        }

        try {
            $obit->update(['is_synchronized' => Obit::SYNCHRONIZED]);
        } catch (Throwable $t) {
            Log::error("Cannot commit gateway obit record", [$t]);

            throw $t;
        }
    }
}
