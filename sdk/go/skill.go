package multiagentspec

// Skill represents a reusable capability that agents can invoke.
// Skills encapsulate instructions, scripts, and resources for specific tasks.
type Skill struct {
	// Name is the unique identifier for the skill (lowercase, hyphenated).
	Name string `json:"name" yaml:"name" jsonschema:"required"`

	// Description is a brief summary of what the skill does.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Instructions is the main prompt/instructions for the skill.
	// This is what gets executed when the skill is invoked.
	Instructions string `json:"instructions,omitempty" yaml:"instructions,omitempty"`

	// Scripts are paths to executable scripts relative to the skill directory.
	// These scripts can be invoked as part of the skill execution.
	Scripts []string `json:"scripts,omitempty" yaml:"scripts,omitempty"`

	// References are paths to reference documentation relative to the skill directory.
	// These provide context and examples for the skill.
	References []string `json:"references,omitempty" yaml:"references,omitempty"`

	// Assets are paths to templates, config files, or other resources.
	// These are files the skill may read or use as templates.
	Assets []string `json:"assets,omitempty" yaml:"assets,omitempty"`

	// Triggers are keywords or patterns that invoke this skill.
	// When an agent sees these triggers, it should consider using this skill.
	Triggers []string `json:"triggers,omitempty" yaml:"triggers,omitempty"`

	// Dependencies are external tools or binaries required by this skill.
	// Examples: "git", "npm", "docker", "go"
	Dependencies []string `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`

	// Model is the recommended model tier for this skill.
	// If not set, the agent's default model is used.
	Model Model `json:"model,omitempty" yaml:"model,omitempty"`

	// Tools are the tools this skill requires access to.
	Tools []string `json:"tools,omitempty" yaml:"tools,omitempty"`
}

// NewSkill creates a new Skill with the given name and description.
func NewSkill(name, description string) *Skill {
	return &Skill{
		Name:        name,
		Description: description,
	}
}

// WithDescription sets the skill's description and returns the skill for chaining.
func (s *Skill) WithDescription(description string) *Skill {
	s.Description = description
	return s
}

// WithInstructions sets the skill's instructions and returns the skill for chaining.
func (s *Skill) WithInstructions(instructions string) *Skill {
	s.Instructions = instructions
	return s
}

// WithModel sets the skill's recommended model and returns the skill for chaining.
func (s *Skill) WithModel(model Model) *Skill {
	s.Model = model
	return s
}

// WithTools sets the skill's required tools and returns the skill for chaining.
func (s *Skill) WithTools(tools ...string) *Skill {
	s.Tools = tools
	return s
}

// AddScript adds a script path to the skill.
func (s *Skill) AddScript(path string) *Skill {
	s.Scripts = append(s.Scripts, path)
	return s
}

// AddReference adds a reference document path to the skill.
func (s *Skill) AddReference(path string) *Skill {
	s.References = append(s.References, path)
	return s
}

// AddAsset adds an asset path to the skill.
func (s *Skill) AddAsset(path string) *Skill {
	s.Assets = append(s.Assets, path)
	return s
}

// AddTrigger adds a trigger keyword to the skill.
func (s *Skill) AddTrigger(keyword string) *Skill {
	s.Triggers = append(s.Triggers, keyword)
	return s
}

// AddDependency adds a dependency to the skill.
func (s *Skill) AddDependency(dep string) *Skill {
	s.Dependencies = append(s.Dependencies, dep)
	return s
}

// HasTrigger returns true if the skill has the given trigger.
func (s *Skill) HasTrigger(trigger string) bool {
	for _, t := range s.Triggers {
		if t == trigger {
			return true
		}
	}
	return false
}

// HasDependency returns true if the skill requires the given dependency.
func (s *Skill) HasDependency(dep string) bool {
	for _, d := range s.Dependencies {
		if d == dep {
			return true
		}
	}
	return false
}

// SkillSet represents a collection of skills that can be loaded and queried.
type SkillSet struct {
	// Skills is a map of skill name to skill definition.
	Skills map[string]*Skill `json:"skills" yaml:"skills"`
}

// NewSkillSet creates a new empty SkillSet.
func NewSkillSet() *SkillSet {
	return &SkillSet{
		Skills: make(map[string]*Skill),
	}
}

// Add adds a skill to the set.
func (ss *SkillSet) Add(skill *Skill) {
	ss.Skills[skill.Name] = skill
}

// Get returns a skill by name, or nil if not found.
func (ss *SkillSet) Get(name string) *Skill {
	return ss.Skills[name]
}

// Names returns all skill names in the set.
func (ss *SkillSet) Names() []string {
	names := make([]string, 0, len(ss.Skills))
	for name := range ss.Skills {
		names = append(names, name)
	}
	return names
}

// FindByTrigger returns all skills that have the given trigger.
func (ss *SkillSet) FindByTrigger(trigger string) []*Skill {
	var matches []*Skill
	for _, skill := range ss.Skills {
		if skill.HasTrigger(trigger) {
			matches = append(matches, skill)
		}
	}
	return matches
}
