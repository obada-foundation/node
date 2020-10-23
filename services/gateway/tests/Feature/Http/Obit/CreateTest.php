<?php

declare(strict_types=1);

namespace Tests\Feature\Http\Obit;

use Laravel\Lumen\Testing\DatabaseTransactions;
use Laravel\Lumen\Testing\WithoutEvents;
use Tests\TestCase;

class CreateTest extends TestCase {

    use DatabaseTransactions, WithoutEvents;

    protected array $validObit = [
        'modified_at'        => '20-03-2020 ',
        'owner_did'          => 'did:obada:owner:123456',
        'serial_number_hash' => 'dc0fb8e9835790195bf4a8e5e122fe608e548f46f88410cc6792927bedbb6d55',
        'manufacturer'       => 'manufacturer',
        'usn'                => '28NwRR9G',
        'part_number'        => 'part number',
        'obit_did'           => '8a53fa0b040e11d71d65554e4a887cad14335f7120345c594af794b2e2808184'
    ];

    /**
     * @test
     */
    public function it_validates_that_obit_has_required_fields() {
        $this->json("POST", route('obits.create'), []);

        $this->seeStatusCode(422);
        $this->seeJson([
            'code'   => 422,
            'errors' => [
                'manufacturer'       => ['The manufacturer field is required.'],
                'modified_at'        => ['The modified at field is required.'],
                'obit_did'           => ['The obit did field is required.'],
                'owner_did'          => ['The owner did field is required.'],
                'part_number'        => ['The part number field is required.'],
                'serial_number_hash' => ['The serial number hash field is required.'],
                'usn'                => ['The usn field is required.'],
            ],
            'message' => null
        ]);
    }

    /**
     * @test
     */
    public function it_validates_modified_at_parameter_is_actual_date() {
        $payload = [
            'obit_did'           => 'did:obada:fe096095-e0f0-4918-9607-6567bd5756b5',
            'manufacturer'       => 'Sony',
            'owner_did'          => 'did:obada:owner:123456',
            'part_number'        => 'MWCN2LL/A',
            'serial_number_hash' => 'f6fc84c9f21c24907d6bee6eec38cabab5fa9a7be8c4a7827fe9e56f245bd2d5',
            'usn'                => '2zEz-xLJR',
            'modified_at'        => 'some date'
        ];

        $this->json("POST", route('obits.create'), $payload);

        $this->seeStatusCode(422);
        $this->seeJson([
            'code'   => 422,
            'errors' => [
                'modified_at' => ['The modified at is not a valid date.'],
            ],
            'message' => null
        ]);
    }

    /**
     * @test
     */
    public function it_returns_correct_response_when_create_basic_obit() {
        $this->json("POST", route('obits.create'), $this->validObit);

        $this->seeStatusCode(204);
    }
}
