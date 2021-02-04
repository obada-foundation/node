<?php

namespace Database\Factories;

use App\Services\Gateway\Models\Obit;
use Illuminate\Database\Eloquent\Factories\Factory;

class ObitFactory extends Factory {

    protected $model = Obit::class;

    public function definition() {
        return [
            'modified_at'        => $this->faker->dateTime(),
            'owner_did'          => 'did:obada:owner:123456',
            'obd_did'            => 'did:obada:owner:123456',
            'serial_number_hash' => 'dc0fb8e9835790195bf4a8e5e122fe608e548f46f88410cc6792927bedbb6d55',
            'manufacturer'       => 'manufacturer',
            'usn'                => '28NwRR9G',
            'part_number'        => 'part number',
            'obit_did'           => '8a53fa0b040e11d71d65554e4a887cad14335f7120345c594af794b2e2808184',
            'root_hash'          => 'ddd'
        ];
    }
}
