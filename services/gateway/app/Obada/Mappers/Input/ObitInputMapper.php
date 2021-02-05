<?php

declare(strict_types=1);

namespace App\Obada\Mappers\Input;

use Carbon\Carbon;
use Illuminate\Support\Arr;
use Illuminate\Support\Facades\Validator;
use Obada\Mappers\Input\InputMapper;
use Obada\Obit;

class ObitInputMapper implements InputMapper {
    /**
     * @param $input
     * @return Obit
     * @throws \Carbon\Exceptions\InvalidFormatException
     * @throws \Illuminate\Validation\ValidationException
     * @throws \Obada\Exceptions\PropertyValidationException
     */
    public function map($input): Obit {
        $obit = Obit::make([
            'manufacturer'       => Arr::get($input, 'manufacturer', ''),
            'serial_number_hash' => Arr::get($input, 'serial_number_hash', ''),
            'part_number'        => Arr::get($input, 'part_number', ''),
            'owner_did'          => Arr::get($input, 'owner_did', ''),
            'obd_did'            => Arr::get($input, 'obd_did', ''),
            'modified_at'        => Carbon::parse(Arr::get($input, 'modified_at')),
            'obit_status'        => Arr::get($input, 'obit_status', ''),
            'metadata'           => Arr::get($input, 'metadata', []),
            'structured_data'    => Arr::get($input, 'structured_data', []),
            'documents'          => Arr::get($input, 'doc_links', []),
        ]);

        Validator::make(
            ['root_hash'    => (string) Arr::get($input, 'root_hash')],
            ['root_hash'    => 'in:' . (string) $obit->rootHash()],
            ['root_hash.in' => 'Integrity error. Expected and actual hashes are not match.']
        )->validate();

        return $obit;
    }
}
