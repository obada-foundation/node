<?php

declare(strict_types=1);

namespace App\Services\Gateway\Validation\Rules;

use Illuminate\Contracts\Validation\Rule;

class DtoArrayKeys implements Rule {
    protected array $keys;

    public function __construct($keys) {
        $this->keys = $keys;
    }
    /**
     * Determine if the validation rule passes.
     *
     * @param string $attribute
     * @param mixed $value
     *
     * @return bool
     */
    public function passes($attribute, $value) {
        return collect($value)
            ->keys()
            ->filter(fn ($key) => !in_array($key, $this->keys))
            ->count() == 0;
    }

    /**
     * Get the validation error message.
     *
     * @return string
     */
    public function message() {
        return sprintf(
            'The attribute :attribute must have only keys: %s.',
            implode(',', $this->keys)
        );
    }
}

