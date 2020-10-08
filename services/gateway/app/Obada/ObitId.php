<?php

declare(strict_types=1);

namespace App\Obada;

use Illuminate\Support\Facades\Validator;
use Tuupola\Base58;

class ObitId {

    protected string $hash;

    /**
     * ObitId constructor.
     * @param string $serialNumberHash
     * @param string $manufacturer
     * @param string $partNumber
     */
    public function __construct(string $serialNumberHash, string $manufacturer, string $partNumber) {
        //$this->validate($serialNumberHash);

        $this->hash = hash('sha256', $manufacturer . $partNumber . $serialNumberHash);
    }

    protected function validate(string $serialNumberHash) {
        Validator::make(
            ['serial_number_hash'       => $serialNumberHash],
            ['serial_number_hash'       => 'required|regex:/^([a-f0-9]{64})$/'],
            ['serial_number_hash.regex' => 'Received the invalid serial number hash. Must be the valid SHA256 hash.']
        )->validate();
    }

    /**
     * Convert Obit to Universal Serial Number
     *
     * @return string
     */
    public function toUsn(): string {
        $encoder = new Base58(["characters" => Base58::BITCOIN]);

        return substr($encoder->encode($this->hash), 0, 8);
    }

    /**
     * @return string
     */
    public function toHash(): string {
        return $this->hash;
    }

    /**
     * @return string
     */
    public function toDid(): string {
        return sprintf('did:obada:%s', $this->hash);
    }
}
