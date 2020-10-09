<?php

declare(strict_types=1);

namespace App\Services\Gateway\Validation\Rules;

use Illuminate\Contracts\Validation\Rule;
use Throwable;

class HashLink implements Rule {
    protected ObitId $obit;

    /**
     * Determine if the validation rule passes.
     *
     * @param string $attribute
     * @param string $url
     *
     * @return bool
     */
    public function passes($attribute, $url) {
        $query = parse_url($url, PHP_URL_QUERY);

        if (! $query) {
            return false;
        }

        parse_str($query, $query);

        if (! isset($query['hl']) || $query['hl'] === null || $query['hl'] === '') {
            return false;
        }

        $documentHash = $query['hl'];

        try {
            $data = file_get_contents($url);
        } catch (Throwable $t) {
            // TODO: log this error

            return false;
        }

        return hash('sha256', $data) === $documentHash;
    }

    /**
     * Get the validation error message.
     *
     * @return string
     */
    public function message() {
        return 'Not a valid hashlink.';
    }
}

