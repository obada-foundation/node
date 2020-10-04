<?php

declare(strict_types=1);

namespace App\Services\Gateway\Validation\Rules;

use Illuminate\Contracts\Validation\Rule;

class HashLink implements Rule {
    /**
     * Determine if the validation rule passes.
     *
     * @param string $attribute
     * @param mixed $value
     *
     * @return bool
     */
    public function passes($attribute, $value) {
        return $value == $this->obit->toHash();
    }

    /**
     * Get the validation error message.
     *
     * @return string
     */
    public function message() {
        return 'Integrity of obit id is broken.';
    }
}

