<?php

declare(strict_types=1);

namespace App\Services\Gateway;

use App\Obada\ObitId;
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
        ]);
    }

    protected function validate() {
        $data  = [

        ];

        $rules = [

        ];

        Validator::make($data, $rules)->validate();
    }
}
