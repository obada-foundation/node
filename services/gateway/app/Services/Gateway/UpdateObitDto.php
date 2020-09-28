<?php

declare(strict_types=1);

namespace App\Services\Gateway;

use App\Services\Gateway\Models\Obit;
use Illuminate\Support\Facades\Validator;
use Spatie\DataTransferObject\DataTransferObject;
use Laravel\Lumen\Http\Request;

class UpdateObitDto extends DataTransferObject {

    public ?string $obitDID;

    public ?string $usn;

    public ?string $modifiedAt;

    public ?string $serialNumberHash;

    public ?string $manufacturer;

    public ?string $partNumber;

    public ?string $obitStatus;

    public ?string $ownerDID;

    public ?string $obdDID;

    public $metadata;

    public $docLinks;

    public $structuredData;

    /**
     * ObitDto constructor.
     * @param array $parameters
     */
    public function __construct(array $parameters = []) {
        parent::__construct($parameters);

        $this->validate();
    }

    /**
     * @param Request $request
     * @return static
     */
    public static function fromRequest(Request $request): self {
        return new self([
            'modifiedAt'       => $request->json('modified_at'),
            'serialNumberHash' => $request->json('serial_number_hash'),
            'manufacturer'     => $request->json('manufacturer'),
            'partNumber'       => $request->json('part_number'),
            'obitStatus'       => $request->json('obit_status'),
            'ownerDID'         => $request->json('owner_did'),
            'obdDID'           => $request->json('obd_did'),
            'metadata'         => $request->json('metadata'),
            'docLinks'         => $request->json('doc_links'),
            'structuredData'   => $request->json('structured_data'),
        ]);
    }

    protected function validate() {
        $data  = [
            'obit_status'        => $this->obitStatus,
            'metadata'           => $this->metadata,
            'metadata'           => $this->metadata,
            'doc_links'          => $this->docLinks,
            'structured_data'    => $this->structuredData
        ];

        $rules = [
            'obit_status'        => 'nullable|in:' . implode(',', Obit::STATUSES),
            'metadata'           => 'nullable|array',
            'doc_links'          => 'nullable|array',
            'structured_data'    => 'nullable|array'
        ];

        Validator::make($data, $rules)->validate();
    }
}
