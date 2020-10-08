<?php

declare(strict_types=1);

namespace App\Services\Gateway;

use App\Obada\ObitId;
use App\Services\Gateway\Validation\Rules\ObitIntegrity;
use Illuminate\Support\Facades\Validator;
use App\Services\Gateway\Models\Obit;
use Laravel\Lumen\Http\Request;

class ObitDto extends BaseDto {

    public string $obitDID = "";

    public string $usn = "";

    public string $modifiedAt = "";

    public string $serialNumberHash = "";

    public string $manufacturer = "";

    public string $partNumber = "";

    public ?string $obitStatus;

    public ?string $ownerDID;

    public ?string $obdDID;

    public ?ObitId $obit;

    /**
     * ObitDto constructor.
     * @param array $parameters
     */
    public function __construct(array $parameters = []) {
        parent::__construct($parameters);

        $this->obit = new ObitId($this->serialNumberHash, $this->manufacturer, $this->partNumber);

        $this->validate();
    }

    /**
     * @param Request $request
     * @return static
     */
    public static function fromRequest(Request $request): self {
        return new self([
            'obitDID'          => $request->json('obit_did'),
            'usn'              => $request->json('usn'),
            'modifiedAt'       => $request->json('modified_at'),
            'serialNumberHash' => $request->json('serial_number_hash'),
            'manufacturer'     => $request->json('manufacturer'),
            'partNumber'       => $request->json('part_number'),
            'obitStatus'       => $request->json('obit_status'),
            'ownerDID'         => $request->json('owner_did'),
            'obdDID'           => $request->json('obd_did'),
            'metadata'         => $request->json('metadata', []),
            'docLinks'         => $request->json('doc_links', []),
            'structuredData'   => $request->json('structured_data', []),
        ]);
    }

    protected function validate() {
        parent::validate();

        $data  = [
            'obit_did'           => $this->obitDID,
            'manufacturer'       => $this->manufacturer,
            'usn'                => $this->usn,
            'modified_at'        => $this->modifiedAt,
            'serial_number_hash' => $this->serialNumberHash,
            'part_number'        => $this->partNumber,
            'obit'               => $this->obit,
            'obit_status'        => $this->obitStatus
        ];

        $rules = [
            'obit_did'           => ['required', new ObitIntegrity($this->obit)],
            'manufacturer'       => 'required',
            'usn'                => 'required|in:' . $this->obit->toUsn(),
            'modified_at'        => 'required',
            'serial_number_hash' => 'required',
            'part_number'        => 'required',
            'obit'               => 'required',
            'obit_status'        => 'nullable|in:' . implode(',', Obit::STATUSES)
        ];

        Validator::make($data, $rules)->validate();
    }

    /**
     * @param string $property
     * @return string
     */
    public function __get(string $property) {
        if (property_exists(self::class, $property)) {
            return $property;
        }

        $field = collect(explode('_', $property))
            ->map(fn($v, $k) => $k !== 0 ? ucfirst($v) : $v)
            ->implode('');

        return $this->$field;
    }
}
