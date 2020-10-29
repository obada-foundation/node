<?php

declare(strict_types=1);

namespace App\Services\Gateway;

use Exception;
use App\Services\Gateway\Events\RecordUpdated;
use App\Services\Gateway\Repositories\GatewayRepositoryContract;
use App\Services\Gateway\Events\RecordCreated;
use App\Services\Gateway\Models\Obit as Model;
use Obada\Obit;
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
     * @param array $args
     * @return \Illuminate\Pagination\LengthAwarePaginator
     */
    public function search(array $args = []) {
        return $this->repository->findBy($args);
    }

    /**
     * @param Obit $o
     * @return Model
     */
    public function create(Obit $o): Model {
        $status = (string) $o->getStatus();

        $obit = new Model;
        $obit->parent_id          = optional(Model::orderBy('id', 'DESC')->first())->id;
        $obit->obit_did           = (string) $o->getObitId()->toDid();
        $obit->usn                = (string) $o->getObitId()->toUsn();
        $obit->obit_status        = $status ?: null;
        $obit->owner_did          = (string) $o->getOwnerDid();
        $obit->obd_did            = (string) $o->getObdDid();
        $obit->metadata           = $o->getMetadata();
        $obit->doc_links          = $o->getDocuments();
        $obit->structured_data    = $o->getStructuredData();
        $obit->manufacturer       = (string) $o->getManufacturer();
        $obit->part_number        = (string) $o->getPartNumber();
        $obit->serial_number_hash = (string) $o->getSerialNumberHash();
        $obit->modified_at        = (string) $o->getModifiedAt();
        $obit->root_hash          = (string) $o->rootHash();
        $obit->save();

        return $obit;
    }

    /**
     * @param string $obitId
     * @param Obit $o
     * @return mixed|void
     */
    public function update(string $obitId, Obit $o) {
        $status = (string) $o->getStatus();

        $obit = $this->repository->find($obitId);

        $obit->obit_status     = $status ?: null;
        $obit->owner_did       = (string) $o->getOwnerDid();
        $obit->obd_did         = (string) $o->getObdDid();
        $obit->metadata        = $o->getMetadata();
        $obit->doc_links       = $o->getDocuments();
        $obit->structured_data = $o->getStructuredData();
        $obit->modified_at     = (string) $o->getModifiedAt();
        $obit->root_hash       = (string) $o->rootHash();
        $obit->is_synchronized = Model::NOT_SYNCHRONIZED;
        $obit->save();
    }

    /**
     * @param string $obitDID
     * @return Model
     */
    public function show(string $obitDID): ?Model {
        return $this->repository->find($obitDID);
    }

    public function delete(string $obitId) {
        $obit = $this->repository->find($obitId);
        $obit->obit_status = Model::DISABLED_BY_OWNER_STATUS;
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
            $obit->update(['is_synchronized' => Model::SYNCHRONIZED]);
        } catch (Throwable $t) {
            Log::error("Cannot commit gateway obit record", [$t]);

            throw $t;
        }
    }
}
