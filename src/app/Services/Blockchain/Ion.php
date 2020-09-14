<?php

declare(strict_types=1);

namespace App\Services\Blockchain;

class Ion {
    /**
     * @var object The go extension that allows to work with AWS Ion
     */
    protected $ion;

    public function __construct() {
        $module = \phpgo_load("/ext/ion.so", "ion");
        $ref = new \ReflectionClass($module);

        $this->ion = $module;
    }

    /**
     * Receive a string of bytes (ion object) and returns PHP array
     *
     * @param string $bytes
     * @return array
     */
    public function decode(string $bytes): array {
        $json = $this->ion->decode($bytes);

        return json_decode($json, true);
    }

    /**
     * @param array $json
     * @return mixed
     */
    public function encode(array $json): string {
        $bytes = $this->ion->encode(json_encode($json));

        return $bytes;
    }
}
